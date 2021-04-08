package utils

import (
	"encoding/hex"
	"github.com/satori/go.uuid"
)

func UUIDFromString(s string) (uuid.UUID, error) {
	if len(s) == len(uuid.Nil.String()) {
		return uuid.FromString(s)
	}

	b, err := hex.DecodeString(s)

	if err != nil {
		return uuid.Nil, err
	}

	return uuid.FromBytes(b)
}