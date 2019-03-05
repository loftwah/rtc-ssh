﻿package main

import (
	"fmt"
	"log"
	"net"
	"io"
 	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
	"github.com/pions/webrtc/pkg/datachannel"
	"github.com/gorilla/websocket"
)

type sendWrap struct {
	*webrtc.RTCDataChannel
}
type errorString struct {
      s string
}
func (e *errorString) Error() string {
      return e.s
}
func (s *sendWrap) Write(b []byte) (int, error) {
	err := s.RTCDataChannel.Send(datachannel.PayloadBinary{Data: b})
	return len(b), err
}
/*
var (
	configRTC = webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
				},
			},
		},
	}
)
*/
func interpreter(c *websocket.Conn, data Json, conf Config) error {
	if data.Error != ""{
		return &errorString{data.Error}
	}
	
	switch data.Type {	
		case "signal_OK":
			log.Println("Signal OK")
		
		case "offer":
			pc, err := webrtc.New(configRTC)
			if err != nil {
				log.Println(err)
				return err
			}
			ssh, err := net.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
			if err != nil {
				log.Println("ssh dial failed:", err)
				pc.Close() 
			}
						
			pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
				log.Println("ICE Connection State has changed:", state)
				if state == ice.ConnectionStateDisconnected {
					pc.Close()
					ssh.Close()
				}
			})
			pc.OnDataChannel(func(dc *webrtc.RTCDataChannel) {
				if dc.Label == "SSH" {
					DataChannel(dc, ssh)
				}
			})
			
			if err := pc.SetRemoteDescription(webrtc.RTCSessionDescription{
				Type: webrtc.RTCSdpTypeOffer,
				Sdp:  data.Sdp,
			}); err != nil {
				log.Println("rtc error:", err)
				pc.Close()
				ssh.Close()
				return err
			}
		
			answer, err := pc.CreateAnswer(nil)
			if err != nil {
				log.Println("rtc error:", err)
				pc.Close()
				ssh.Close()
				return err
			}
			if err = c.WriteJSON(answer); err != nil {
				log.Println("write signal error:", err)
				ssh.Close()
				return err
			}
		default:
			return &errorString{"Not interpreted"}
	}
	return nil
}


func DataChannel(dc *webrtc.RTCDataChannel, ssh net.Conn) {
		dc.OnOpen(func() {	
			log.Println("Connect SSH socket")
			message := "OPEN_RTC_CHANNEL"
			err := dc.Send(datachannel.PayloadString{Data: []byte(message)})
			if err != nil{
				log.Println("write data error:", err)
			}
			io.Copy(&sendWrap{dc}, ssh)
		})
		
		dc.OnMessage(func(payload datachannel.Payload) {
			switch p := payload.(type) {
				case *datachannel.PayloadString:				
					log.Printf("\nReceive: %s\n", string(p.Data))
				case *datachannel.PayloadBinary:
					_, err := ssh.Write(p.Data)
					if err != nil {
						log.Println("ssh write failed:", err)
						return
					}
				default:
					log.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), dc.Label)
			}
		})
		/*
		dc.OnClose(func() {
			log.Printf("Close data channel '%s' ID: %d\n", dc.Label, dc.ID)
			ssh.Close()
		})
		*/
}
