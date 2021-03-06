package bip39go

import (
	"encoding/hex"
	"errors"
	"fmt"
)

type ConfidentialData interface {
	Free()
}

type HexConvertible interface {
	// ToHexString - convert raw bytes to hex string
	ToHexString() string
	FromHexString(hex string) error
}

// MinterData64 - 64 bytes Data array
type MinterData64 struct{ Data [64]uint8 }

// MinterData33 - 33 bytes public key array
type MinterData33 struct{ Data [33]uint8 }

// MinterData32 - 32 bytes common array
type MinterData32 struct{ Data [32]uint8 }

// MinterBip32Key - 112 bytes special array type
type MinterBip32Key struct{ Data [112]uint8 }

func (target *MinterData64) ToHexString() string {
	return hex.EncodeToString(target.Data[:])
}

func (target *MinterData64) FromHexString(h string) error {
	return fromHexWithSize(target.Data[:], h, 64)
}

func (target *MinterData64) Free() {
	memset(target.Data[:], 0x00)
}

func (target *MinterData33) ToHexString() string {
	return hex.EncodeToString(target.Data[:])
}

func (target *MinterData33) FromHexString(h string) error {
	return fromHexWithSize(target.Data[:], h, 33)
}

func (target *MinterData33) Free() {
	memset(target.Data[:], 0x00)
}

func (target *MinterData32) ToHexString() string {
	return hex.EncodeToString(target.Data[:])
}

func (target *MinterData32) FromHexString(h string) error {
	return fromHexWithSize(target.Data[:], h, 32)
}

func (target *MinterData32) Free() {
	memset(target.Data[:], 0x00)
}

func (target *MinterBip32Key) ToHexString() string {
	return hex.EncodeToString(target.Data[:])
}

func (target *MinterBip32Key) FromHexString(h string) error {
	return fromHexWithSize(target.Data[:], h, 112)
}

func (target *MinterBip32Key) Free() {
	memset(target.Data[:], 0x00)
}

func fromHexWithSize(target []uint8, h string, size int) error {
	if len(h) != (size * 2) {
		return errors.New(fmt.Sprintf("Invalid input length: must be %d chars (%d bytes), given: %d", size*2, size, len(h)))
	}

	res, err := hex.DecodeString(h)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to decode hex string: %v", err))
	}

	for idx, val := range res {
		target[idx] = uint8(val)
	}

	return nil
}
