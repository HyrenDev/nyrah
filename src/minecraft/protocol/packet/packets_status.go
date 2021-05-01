package packet

import (
	chat2 "net/hyren/nyrah/minecraft/chat"
	codecs2 "net/hyren/nyrah/minecraft/protocol/codecs"
)

type StatusRequest struct{}

func (_ StatusRequest) ID() int { return 0x00 }

type StatusResponse struct {
	Status struct {
		Version struct {
			Name     string `json:"name"`
			Protocol int    `json:"protocol"`
		} `json:"version"`

		Players struct {
			Max    int `json:"max"`
			Online int `json:"online"`
		} `json:"players"`

		ModInfo struct {
			Type    string   `json:"type"`
			ModList []string `json:"modList"`
		} `json:"modinfo"`

		Favicon string `json:"favicon"`

		Description chat2.TextComponent `json:"description"`
	}
}

func (_ StatusResponse) ID() int { return 0x00 }

type StatusPing struct {
	Payload codecs2.Long
}

func (_ StatusPing) ID() int { return 0x01 }

type StatusPong struct {
	Payload codecs2.Long
}

func (_ StatusPong) ID() int { return 0x01 }
