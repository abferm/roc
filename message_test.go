package roc

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestMessage_Validate(t *testing.T) {
	m := Message{}
	m.Destination = Address{Unit: 1, Group: 2}
	m.Source = Address{Unit: 1, Group: 0}
	m.Opcode = 17
	m.DataLength = 3
	m.Data = []byte{'M', 'O', 'C'}
	m.CRC = 0x1885

	assert.NoError(t, m.validate())
}

func TestMessage_Validate_BadLength(t *testing.T) {
	m := Message{}
	m.Destination = Address{Unit: 1, Group: 2}
	m.Source = Address{Unit: 1, Group: 0}
	m.Opcode = 17
	m.DataLength = 3
	m.Data = []byte{'M', 'O', 'C', 0}
	m.CRC = 0x1885

	if assert.Error(t, m.validate()) {
		assert.EqualError(t, m.validate(), "Data Length Mismatch")
	}
}

func TestMessage_Read(t *testing.T) {
	expected := Message{}
	expected.Destination = Address{Unit: 1, Group: 2}
	expected.Source = Address{Unit: 1, Group: 0}
	expected.Opcode = 17
	expected.DataLength = 3
	expected.Data = []byte{'M', 'O', 'C'}
	expected.CRC = 0x1885

	data := []byte{1, 2, 1, 0, 17, 3, 'M', 'O', 'C', 0x85, 0x18}
	m := Message{}

	assert.NoError(t, m.read(bytes.NewBuffer(data)))

	assert.EqualValues(t, expected, m)
}

func TestMessage_Read_BadCRC(t *testing.T) {
	expected := Message{}
	expected.Destination = Address{Unit: 1, Group: 2}
	expected.Source = Address{Unit: 1, Group: 0}
	expected.Opcode = 17
	expected.DataLength = 3
	expected.Data = []byte{'M', 'O', 'C'}
	expected.CRC = 0x1886

	data := []byte{1, 2, 1, 0, 17, 3, 'M', 'O', 'C', 0x86, 0x18}
	m := Message{}

	err := m.read(bytes.NewBuffer(data))

	if assert.Error(t, err) {
		assert.EqualError(t, err, "CRC Mismatch")
	}

	assert.EqualValues(t, expected, m)
}
