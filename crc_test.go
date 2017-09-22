// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package roc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectByte(t *testing.T) {
	testInputs := []byte{0x01, 0xF0, 0x0F, 0xFF}
	testOutputs := []byte{0x80, 0x0F, 0xF0, 0xFF}
	for i, b := range testInputs {
		assert.Equal(t, testOutputs[i], reflectByte(b))
	}
}

func TestReflectUint16(t *testing.T) {
	testInputs := []uint16{0x0001, 0xF0F0, 0x0F0F, 0xFF00}
	testOutputs := []uint16{0x8000, 0x0F0F, 0xF0F0, 0x00FF}
	for i, b := range testInputs {
		assert.Equal(t, testOutputs[i], reflectUint16(b))
	}
}

func TestCRC(t *testing.T) {
	const expected uint16 = 0x1885

	value := calculateCRC([]byte{1, 2, 1, 0, 17, 3, 'M', 'O', 'C'})

	if expected != value {
		t.Fatalf("crc expected 0x%x, actual 0x%x", expected, value)
	}
}
