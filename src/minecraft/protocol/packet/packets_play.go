package packet

import (
	chat2 "net/hyren/nyrah/minecraft/chat"
	codecs2 "net/hyren/nyrah/minecraft/protocol/codecs"
)

type PlayKeepAlive struct {
	AliveId codecs2.VarInt
}

func (_ PlayKeepAlive) ID() int { return 0x1F }

type PlayChatMessage struct {
	Chat     chat2.TextComponent
	Position codecs2.Byte
}

func (_ PlayChatMessage) ID() int { return 0x0F }

type PlayJoinGame struct {
	EntityId   codecs2.Int
	Gamemode   codecs2.UnsignedByte
	Dimension  codecs2.Int
	Difficulty codecs2.UnsignedByte
	MaxPlayers codecs2.UnsignedByte
	LevelType  codecs2.String
	Debug      codecs2.Boolean
}

func (_ PlayJoinGame) ID() int { return 0x23 }

type PlaySpawnPosition struct {
	Location codecs2.Long
}

func (_ PlaySpawnPosition) ID() int { return 0x43 }

type PlayPositionAndLook struct {
	X     codecs2.Double
	Y     codecs2.Double
	Z     codecs2.Double
	Yaw   codecs2.Float
	Pitch codecs2.Float
	Flags codecs2.Byte
	Data  codecs2.VarInt
}

func (_ PlayPositionAndLook) ID() int { return 0x2E }
