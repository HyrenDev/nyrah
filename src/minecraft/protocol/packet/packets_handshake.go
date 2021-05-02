package packet

import (
	codecs2 "net/hyren/nyrah/minecraft/protocol/codecs"
)

type Handshake struct {
	ProtocolVersion codecs2.VarInt
	ServerAddress   codecs2.String
	ServerPort      codecs2.UnsignedShort
	NextState       codecs2.VarInt
}

func (_ Handshake) ID() int { return 0x00 }
