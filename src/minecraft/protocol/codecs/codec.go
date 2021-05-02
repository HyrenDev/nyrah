package codecs

import (
	"errors"
	"io"
)

var UnknownCodecType = errors.New("unknown codec type")

type Codec interface {
	Decode(r io.Reader) (interface{}, error)
	Encode(w io.Writer) error
}
