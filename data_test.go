package bip39go

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMinterData64_ToHexString(t *testing.T) {
	var randBytes = make([]uint8, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	data := MinterData64{}
	copy(data.Data[:], randBytes[:])

	src := data.ToHexString()

	firstByteString := src[0:2]
	firstByte, err := strconv.ParseInt(firstByteString, 16, 64)

	assert.Nil(t, err)
	assert.Equal(t, randBytes[0], uint8(firstByte))

	lastByteString := src[len(src)-2:]
	lastByte, err := strconv.ParseInt(lastByteString, 16, 64)
	assert.Nil(t, err)
	assert.Equal(t, randBytes[len(randBytes)-1], uint8(lastByte))

	data.Free()
	assert.Equal(t, data.Data[0], uint8(0x00))
	assert.Equal(t, data.Data[len(data.Data)-1], uint8(0x00))
}

func TestMinterData32_ToHexStringEmpty(t *testing.T) {
	data := MinterData32{}
	assert.Equal(t, 32, len(data.Data))
	hex := data.ToHexString()
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", hex)

	data.Free()
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", hex)
}

func TestMinterData64_FromHexString(t *testing.T) {
	src := "fd90261f5bd702ffbe7483c3b5aa7b76b1f40c1582cc6a598120b16067d3cb9afd90261f5bd702ffbe7483c3b5aa7b76b1f40c1582cc6a598120b16067d3cb9a"
	data := MinterData64{}
	err := data.FromHexString(src)

	assert.Nil(t, err)

	firstByteString := hex.EncodeToString(data.Data[0:1])
	assert.Equal(t, firstByteString, src[0:2])

	lastByteString := hex.EncodeToString(data.Data[len(data.Data)-1:])
	assert.Equal(t, lastByteString, src[len(src)-2:])

	data.Free()
	assert.Equal(t, data.Data[0], uint8(0x00))
	assert.Equal(t, data.Data[len(data.Data)-1], uint8(0x00))
}

func TestMinterData64_FromHexStringInvalidLength(t *testing.T) {
	src := "fd90261f5bd702ffbe7483c3b5aa7b76b1f40c1582cc6a598120b16067d3cb9a"
	data := MinterData64{}
	err := data.FromHexString(src)

	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestMinterData64_FromHexStringIncompatibleChars(t *testing.T) {
	src := "fd90261f5bd702__be7483c3b5aa7b$$b1f40c1582cc6a598120b16067d3cb9afd90261f5bd702__be7483c3b5aa7b$$b1f40c1582cc6a598120b16067d3cb9a"
	data := MinterData64{}
	err := data.FromHexString(src)

	assert.NotNil(t, err)
	fmt.Println(err)
}
