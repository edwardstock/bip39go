package bip39go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemset(t *testing.T) {
	data := MinterData64{}
	data.Data[0] = 0xff
	data.Data[1] = 0xff - 1
	data.Data[2] = 0xff - 2
	data.Data[3] = 0xff - 3

	assert.Equal(t, data.Data[0], uint8(0xff))
	assert.Equal(t, data.Data[1], uint8(0xff-1))
	assert.Equal(t, data.Data[2], uint8(0xff-2))
	assert.Equal(t, data.Data[3], uint8(0xff-3))

	memset(data.Data[:], 0x00)

	assert.Equal(t, data.Data[0], uint8(0x00))
	assert.Equal(t, data.Data[1], uint8(0x00))
	assert.Equal(t, data.Data[2], uint8(0x00))
	assert.Equal(t, data.Data[3], uint8(0x00))

	another := make([]uint8, 0)
	memset(another, 0x00)
}

func TestMemsetString(t *testing.T) {

	somehtingString := make([]string, 3)
	somehtingString[0] = "aaa"
	somehtingString[1] = "bbb"
	somehtingString[2] = "ccc"

	memsetString(somehtingString[:], "")

	assert.Equal(t, somehtingString[0], "")
	assert.Equal(t, somehtingString[1], "")
	assert.Equal(t, somehtingString[2], "")

	another := make([]string, 0)
	memsetString(another, "")
}
