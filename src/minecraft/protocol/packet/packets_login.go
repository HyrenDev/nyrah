package packet

import (
	"../../chat"
	"../../protocol/codecs"
)

type LoginStart struct {
	Username codecs.String
}

func (_ LoginStart) ID() int { return 0x00 }

type LoginSuccess struct {
	UUID     codecs.String
	Username codecs.String
}

func (_ LoginSuccess) ID() int { return 0x02 }

type LoginDisconnect struct {
	Chat chat.TextComponent
}

func (_ LoginDisconnect) ID() int { return 0x00 }
