package roc

import (
	"encoding/binary"
	"errors"
	"io"
)

const MaxPayloadSize = 0xffff

type header struct {
	Destination Address
	Source      Address
	Opcode      byte
	DataLength  byte
}

func (header *header) read(reader io.Reader) error {
	// Read destination address
	err := header.Destination.read(reader)
	if err != nil {
		return err
	}

	// Read source address
	err = header.Source.read(reader)
	if err != nil {
		return err
	}

	// read Opcode and DataLength
	remainingHeader := make([]byte, 2)
	_, err = io.ReadFull(reader, remainingHeader)
	header.Opcode, header.DataLength = remainingHeader[0], remainingHeader[1]
	return err
}

func (header *header) bytes() []byte {
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
	err := message.header.read(reader)
	if err != nil {
		return err
	}

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

func (message *Message) bytes() []byte {
	data := append(message.header.bytes(), message.Data...)
	crcBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcBytes, message.CRC)
	return append(data, crcBytes...)
}

func (message *Message) validate() (err error) {
	if len(message.Data) != int(message.DataLength) {
		return errors.New("Data Length Mismatch")
	}

	data := message.bytes()

	// Calculate CRC for the entire message minus the crc
	crc := calculateCRC(data[:len(data)-2])
	if crc != message.CRC {
		err = errors.New("CRC Mismatch")
	}
	return err
}

func (message *Message) updateCalculatedFields() (err error) {
	if len(message.Data) > MaxPayloadSize {
		return errors.New("Message payload too long")
	}
	message.DataLength = byte(len(message.Data))
	data := message.bytes()

	// Calculate CRC for the entire message minus the crc
	message.CRC = calculateCRC(data[:len(data)-2])
	return
}
