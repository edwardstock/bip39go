package bip39go

type Entropy = int

// Entropy variations, default 128 bit == 16 bytes
const (
	Entropy128bits Entropy = 16
	Entropy160bits         = 20
	Entropy192bits         = 24
	Entropy224bits         = 28
	Entropy256bits         = 32
	Entropy288bits         = 36
	Entropy320bits         = 40
)
