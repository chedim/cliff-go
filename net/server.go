package net;

func init() {
  bin, bout := OpenBroadcastChannel()
  parsedIn := ParsePackets(bin)

  myIp := GetLocalIP()

  for packet := range parsedIn {
    rsp := packet.Message.Apply()
  }
}
