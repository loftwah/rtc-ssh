# Rtc-SSH
SSH via WebRTC (pions) tunnel
Rtc-SSH enables connection with SSH protocol eg Raspberry PI BeagleBone and other devices, servers from the browser using WebRTC. Solves the problem of the lack of public IP address, proxy server, servers behind NAT etc. You can connect to an SSH session: https://www.sqs.io
### Usage

```
$ go get -u github.com/mxseba/rtc-ssh
$ cd $GOPATH/bin
```
### First run
```
$ ./rtc-ssh newkey
uuid xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Signal OK
```
Copy uuid key and paste on the static page: https://www.sqs.io 
