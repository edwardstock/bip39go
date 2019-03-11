package main

import (
	"crypto/rand"
	"fmt"
	"github.com/edwardstock/bip39go"
)

func main() {
	rbytes := make([]uint8, 16)
	_, err := rand.Read(rbytes)
	if err != nil {
		panic(err)
	}

	m, err := bip39go.GenerateMnemonicFromBytes(rbytes, "en", bip39go.Entropy128bits)
	defer m.Free()

	if err != nil {
		panic(err)
	}

	fmt.Println(m.Status)
	fmt.Println(m.Raw)

}
