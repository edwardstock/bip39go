package bip39go

// #cgo darwin LDFLAGS: -L${SRCDIR}/libs/darwin_amd64/lib
// #cgo linux amd64 LDFLAGS: -L${SRCDIR}/libs/linux_amd64/lib
// #cgo darwin CFLAGS: -I${SRCDIR}/libs/darwin_amd64/include
// #cgo linux amd64 CFLAGS: -I${SRCDIR}/libs/linux_amd64/include
// #cgo CFLAGS: -std=c11
// #cgo LDFLAGS: -lbip39_go -ltrezor_crypto -lstdc++
// #include <bip39.h>
// #include <hdkey_encoder.h>
//
import "C"
import (
	"errors"
	"fmt"
)

// CStringArray is a just binding for C-native array.
// Represents native go string
type CStringArray struct {
	Items []string
	ptr   **C.char
	size  C.size_t
}

// Free C-array.
// Use: defer arr.Free()
func (target CStringArray) Free() {
	C.minter_free_string_array(target.ptr, target.size)
	target.Items = nil
}

// Length of C-array
func (target CStringArray) Length() int {
	return int(target.size)
}

// HDKey contains all required private Data
type HDKey struct {
	ptr           *C.minter_hdkey
	PublicKey     MinterData33
	PrivateKey    MinterData32
	ChainCode     MinterData32
	ExtPrivateKey MinterBip32Key
	ExtPublicKey  MinterBip32Key
	Depth         uint8
	Index         uint32
	Fingerprint   uint32
}

// Free - completely frees all confidential Data, including memory cleaning (set zeroes to arrays)
func (target *HDKey) Free() {
	C.free_hdkey(target.ptr)
	memset(target.PublicKey.Data[:], 0x00)
	memset(target.PrivateKey.Data[:], 0x00)
	memset(target.ChainCode.Data[:], 0x00)
	memset(target.ExtPrivateKey.Data[:], 0x00)
	memset(target.ExtPublicKey.Data[:], 0x00)
	target.Depth = 0x00
	target.Index = 0
	target.Fingerprint = 0
}

func MakeBip32RootKey(seed MinterData64) *HDKey {
	seedData := C.minter_data64{}

	for idx, val := range seed.Data {
		seedData.data[idx] = C.uchar(val)
	}

	var res *C.minter_hdkey = C.encoder_make_bip32_root_key(&seedData)
	out := &HDKey{ptr: res}
	fromCArrayByte(res.public_key.data[:], out.PublicKey.Data[:])
	fromCArrayByte(res.private_key.data[:], out.PrivateKey.Data[:])
	fromCArrayByte(res.chain_code.data[:], out.ChainCode.Data[:])
	fromCArrayByte(res.ext_private_key.data[:], out.ExtPrivateKey.Data[:])
	fromCArrayByte(res.ext_public_key.data[:], out.ExtPublicKey.Data[:])

	out.Depth = uint8(res.index)
	out.Index = uint32(res.index)
	out.Fingerprint = uint32(res.fingerprint)

	return out
}

// MakeExtenderKey from derivation path and root derived key
// signaure: minter_hdkey *encoder_make_ext_key(const struct minter_hdkey *root_key, const char *derivation_path);
func MakeExtendedKey(key *HDKey, derivationPath string) *HDKey {

	var out C.minter_hdkey

	var pubKeyData C.minter_data33
	var privKeyData C.minter_data32
	var chainCodeData C.minter_data32
	var extPubKeyData C.minter_bip32_key
	var extPrivKeyData C.minter_bip32_key

	// pub key
	for i := 0; i < len(key.PublicKey.Data); i++ {
		pubKeyData.data[i] = (C.uchar)(key.PublicKey.Data[i])
	}
	out.public_key = pubKeyData

	// priv key
	for i := 0; i < len(key.PrivateKey.Data); i++ {
		privKeyData.data[i] = (C.uchar)(key.PrivateKey.Data[i])
	}
	out.private_key = privKeyData

	// chain code
	for i := 0; i < len(key.ChainCode.Data); i++ {
		chainCodeData.data[i] = (C.uchar)(key.ChainCode.Data[i])
	}
	out.chain_code = chainCodeData

	// ext pub key
	for i := 0; i < len(key.ExtPublicKey.Data); i++ {
		extPubKeyData.data[i] = (C.uchar)(key.ExtPublicKey.Data[i])
	}
	out.ext_public_key = extPubKeyData

	// ext priv key
	for i := 0; i < len(key.ExtPrivateKey.Data); i++ {
		extPrivKeyData.data[i] = (C.uchar)(key.ExtPrivateKey.Data[i])
	}
	out.ext_private_key = extPrivKeyData

	out.depth = (C.uchar)(key.Depth)
	out.index = (C.uint)(key.Index)
	out.fingerprint = (C.uint)(key.Fingerprint)

	var res *C.minter_hdkey = C.encoder_make_ext_key(&out, C.CString(derivationPath))

	extOut := &HDKey{ptr: res}
	fromCArrayByte(res.public_key.data[:], extOut.PublicKey.Data[:])
	fromCArrayByte(res.private_key.data[:], extOut.PrivateKey.Data[:])
	fromCArrayByte(res.chain_code.data[:], extOut.ChainCode.Data[:])
	fromCArrayByte(res.ext_private_key.data[:], extOut.ExtPrivateKey.Data[:])
	fromCArrayByte(res.ext_public_key.data[:], extOut.ExtPublicKey.Data[:])

	extOut.Depth = uint8(res.index)
	extOut.Index = uint32(res.index)
	extOut.Fingerprint = uint32(res.fingerprint)

	return extOut
}

// GetLanguages return list of available bip39 languages
// In current version: 7 languages, including:
//  - English
//  - Spanish
//  - French
//  - Italian
//  - Japanese
//  - Chinese (zhs)
//  - Chinese (zht)
func GetLanguages() CStringArray {
	out := CStringArray{}
	out.ptr = C.minter_get_languages(&out.size)
	out.Items = fromCStringArray(out.size, out.ptr)

	return out
}

// GetWordsFromLanguage return list of words for specified language
func GetWordsFromLanguage(lang string) CStringArray {
	out := CStringArray{}
	out.ptr = C.minter_get_words_from_language(C.CString(lang), &out.size)
	out.Items = fromCStringArray(out.size, out.ptr)

	return out
}

// MnemonicResult result of generated mnemonic
type MnemonicResult struct {
	ptr    *C.minter_mnemonic_result
	Words  []string
	Length int
	Raw    string
	Status MnemonicStatus
}

func (target *MnemonicResult) Free() {
	if target == nil {
		return
	}

	if target.ptr != nil {
		C.minter_free_mnemonic(target.ptr)
	}

	memsetString(target.Words, "")
	target.Length = 0
	target.Status = Ok
	target.Raw = ""
}

// GenerateMnemonicRandom internally uses PCGRandom library
// see http://www.pcg-random.org/
// todo: more accurate api, unified error handling
func GenerateMnemonicRandom(language string, entropy Entropy) (result *MnemonicResult, err error) {

	if !ValidateLanguage(language) {
		return nil, errors.New(fmt.Sprintf("invalid language: %s", language))
	}

	var res *C.minter_mnemonic_result = C.minter_generate_mnemonic(C.CString(language), C.size_t(entropy))
	out := &MnemonicResult{ptr: res}

	out.Status = MnemonicStatus(res.status)
	if out.Status != Ok {
		return nil, errors.New(fmt.Sprintf("unable to generate mnemonic: %s", out.Status.String()))
	}

	out.Length = int(res.len)
	out.Raw = C.GoString(res.raw)
	out.Words = fromCStringArray(res.len, res.words)

	return out, nil
}

func ValidateLanguage(lang string) bool {
	langs := GetLanguages()

	for _, item := range langs.Items {
		if item == lang {
			return true
		}
	}

	return false
}

// GenerateMnemonicFromBytes uses bytes obtained by yourself to generate mnemonic.
//  - bytes []uint8 - length must be equally to entropy.
// Example: If you use Entropy128bits (16 byte) - bytes variable must contains exactly 16 bytes
func GenerateMnemonicFromBytes(bytes []uint8, language string, entropy Entropy) (result *MnemonicResult, err error) {
	if len(bytes) != entropy {
		return nil, errors.New(fmt.Sprintf("Invalid bytes size: must be equals to entropy - %d", entropy))
	}

	// todo: move to utils

	// allocate bytes: uint8_t *arr = (uint8_t*) malloc(sizeof(uint8_t) * len)
	inBytes := NewCByteArray(uint64(len(bytes)))
	// fill C-array with input bytes
	inBytes.Write(bytes)
	// free on exit free(arr)
	defer inBytes.Free()

	var res *C.minter_mnemonic_result = C.minter_encode_bytes(inBytes, C.CString(language), C.size_t(entropy))
	out := &MnemonicResult{ptr: res}
	out.Length = int(res.len)
	out.Raw = C.GoString(res.raw)
	out.Words = fromCStringArray(res.len, res.words)

	return out, nil
}

// WordsToSeed convert mnemonic words to encoded seed.
// This seed using everywhere
func WordsToSeed(words string, data64 *MinterData64, numWritten *int) {
	inArr := NewCByteArray(64)

	//noinspection ALL
	var written C.size_t = 0
	C.minter_words_to_seed(C.CString(words), inArr, &written)
	*numWritten = int(written)

	fromCPtrArray(inArr, data64.Data[:])
}

// ValidateMnemonicWords checks passed words for correctness with internal table. Table contains exactly 2048 words
func ValidateMnemonicWords(language, words string) bool {
	var res C.bool = C.minter_validate_words(C.CString(language), C.CString(words))

	return (bool)(res)
}
