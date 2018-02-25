# Netpipes
Netpipes is a simple tool written in Go for redirecting network traffic.

For now it supports listening to TCP or UDP connections and redirecting them to a specified target.

It can be used as a simple proxy or to redirect IPv4 connections to an IPv6 server or vice versa.

## Usage
Netpipes is a command line tool, requires Go 1.9 to compile. It uses TCP by default but that can be overridden.
```
Usage of netpipes:
  -from string
    	Address to listen on
  -to string
    	Address to redirect incoming connections to
  -udp
    	Create UDP tunnel instead of TCP tunnel
```

Example invocations:
 - `netpipes -from localhost:8000 -to localhost:8080` - makes a local TCP server listening at 8080 also visible at 8000
 - `netpipes -from :8000 -to localhost:8080` - same as above but listens on 8000 on all interfaces (not only loopback)
 - `netpipes -from :8000 -to google.com:80` - redirects requests to our machine at port 8000 to google.com
 - `netpipes -from :8000 -to localhost:8080 -udp` - makes a local UDP server listening at 8080 also visible at 8000

## Ideas for improvement
#### UDP over TCP
Sometimes it's not possible to use UDP because of firewalls, or because you want to go through an SSH tunnel.

A possible next step would be tunneling UDP over TCP - this would require 2 instances of the program running, one listening for UDP packets and sending them encoded over TCP to the other end that would decode them and send as UDP packets to the destination.

#### Encryption
One netpipes process can listen to connections and relay them to another instance through an encrypted connection which then can relay them to the destination.

This would be similar to SSH tunneling.