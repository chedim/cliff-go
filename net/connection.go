package net

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

const broadcast_addr = "255.255.255.255"
const port = "42424"

func OpenBroadcastChannel() (<-chan Packet, chan<- Packet) {
	receive := make(chan Packet, 10)
	send := make(chan Packet, 10)

	go listen(receive)
	go broadcast(send)

	return receive, send
}

func listen(receive chan Packet) {
  localAddress, _ := net.ResolveUDPAddr("udp", port)
  if connection, e := net.ListenUDP("udp", localAddress); e != nil {
    panic(e)
  } else {
    defer connection.Close()
    var message Packet
    for {
      inputBytes := make([]byte, 4096)
      length, _, _ := connection.ReadFromUDP(inputBytes)
      buffer := bytes.NewBuffer(inputBytes[:length])
      decoder := gob.NewDecoder(buffer)
      decoder.Decode(&message)

      receive <- message
    }
  }
}

func broadcast(send chan Packet) {
  broadcastAddress, _ := net.ResolveUDPAddr("udp", broadcast_addr+":"+port)
  localAddress, _ := net.ResolveUDPAddr("udp", GetLocalIP())
  if connection, err := net.DialUDP("udp", localAddress, broadcastAddress); err != nil {
    panic(err)
  } else {
    defer connection.Close()

    var buffer bytes.Buffer
    encoder := gob.NewEncoder(&buffer)
    for {
      message := <-send
      encoder.Encode(message)
      connection.Write(buffer.Bytes())
      buffer.Reset()
    }
  }
}

func GetLocalIP() string {
	var localIP string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("GetLocalIP in communication failed")
		return "localhost"
	}
	for _, val := range addr {
		if ip, ok := val.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				localIP = ip.IP.String()
			}
		}
	}
	return localIP
}
