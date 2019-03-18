package bip39go

// MnemonicStatus gives the status after getting words or something else
type MnemonicStatus int

const (
	Ok MnemonicStatus = iota
	UnsupportedEntropy
	UnknownError
)

func (t MnemonicStatus) String() string {
	switch t {
	case Ok:
		return "Ok"
	case UnsupportedEntropy:
		return "Unsupported Entropy"
	case UnknownError:
		fallthrough
	default:
		return "Unknown Error"
	}
}
