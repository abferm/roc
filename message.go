package roc

import (
	"encoding/binary"
	"errors"
	"io"
)

type header struct {
	Destination Address
	Source      Address
	Opcode      byte
	DataLength  byte
}

func (header header) bytes() []byte {
	addresses := append(header.Destination.bytes(), header.Source.bytes()...)
	return append(addresses, header.Opcode, header.DataLength)
}

type Message struct {
	header
	Data []byte
	CRC  uint16
}

// read reads a message from the provided reader
func (message *Message) read(reader io.Reader) error {
	// Read destination address
	err := message.Destination.read(reader)
	if err != nil {
		return err
	}

	// Read source address
	err = message.Source.read(reader)
	if err != nil {
		return err
	}

	// read Opcode and DataLength
	remainingHeader := make([]byte, 2)
	_, err = io.ReadFull(reader, remainingHeader)
	if err != nil {
		return err
	}
	message.Opcode, message.DataLength = remainingHeader[0], remainingHeader[1]

	// read data
	message.Data = make([]byte, int(message.DataLength))
	_, err = io.ReadFull(reader, message.Data)
	if err != nil {
		return err
	}

	err = binary.Read(reader, binary.LittleEndian, &message.CRC)
	if err != nil {
		return err
	}
	return message.validate()
}

func (message *Message) validate() (err error) {
	if len(message.Data) != int(message.DataLength) {
		return errors.New("Data Length Mismatch")
	}
	data := []byte{message.Destination.Unit, message.Destination.Group, message.Source.Unit, message.Source.Group, message.Opcode, message.DataLength}
	data = append(data, message.Data...)

	crc := calculateCRC(data)
	if crc != message.CRC {
		err = errors.New("CRC Mismatch")
	}
	return err
}
