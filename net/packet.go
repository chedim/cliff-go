package net

import (
	"cliff/core"
	"net"
)

type Packet struct {
	SourceIp net.Addr
	Message  string
}

type ParsedPacket struct {
	Packet
	Message *FFIMessage
}

type FFIMessage interface {
	Apply()
}

func ParsePackets(packets <-chan Packet) <-chan ParsedPacket {
	out := make(chan ParsedPacket)
  go parse(packets, out)
  return out
}

func parse(packets <-chan Packet, out chan<- ParsedPacket) {

	for packet := range packets {
    var parsed = ParsedPacket{
      packet,
      nil,
    }

		if len(packet.Message) < 2 {
      core.LogDebug("Discarding invalid command: %s", packet.Message)
		} else if marker := packet.Message[0]; marker == '?' {
      parsed.Message = ReadQueryMessage(packet.Message)
    } else if marker == '=' {
      parsed.Message = ReadValueMessage(packet.Message)
    } else if marker == '@' {
      parsed.Message = ReadSubscribeMessage(packet.Message)
    } else if marker == '#' {
      parsed.Message = ReadUnsubscribeMessage(packet.Message)
    } else if marker == '-' {
      parsed.Message = ReadUnsetMessage(packet.Message)
    }

    out <- parsed
	}
}
