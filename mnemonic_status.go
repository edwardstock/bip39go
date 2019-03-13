package bip39go

// MnemonicStatus gives the status after getting words or something else
type MnemonicStatus = int

const (
	Ok MnemonicStatus = iota
	UnsupportedEntropy
	UnknownError
)
