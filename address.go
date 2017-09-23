package roc

import (
	"io"
)

const (
	DefaultGroup  byte = 240
	DefaultUnit   byte = 240
	BroadcastUnit byte = 0
)

type Address struct {
	Unit  byte
	Group byte
}

// read sets address fields from the given io.Reader
func (addr *Address) read(reader io.Reader) error {
	data := make([]byte, 2)
	_, err := io.ReadFull(reader, data)
	if err != nil {
		return err
	}
	addr.Unit, addr.Group = data[0], data[1]
	return nil
}

func (addr *Address) bytes() []byte {
	return []byte{addr.Unit, addr.Group}
}
