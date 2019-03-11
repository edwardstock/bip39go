package bip39go

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLanguages(t *testing.T) {
	langs := GetLanguages()
	defer langs.Free()

	for _, val := range langs.Items {
		fmt.Println(val)
	}

	assert.Equal(t, 7, len(langs.Items))
	assert.Equal(t, English, langs.Items[0])
	assert.Equal(t, ChineseZHT, langs.Items[6])
}

func TestGetWordsFromLanguage(t *testing.T) {
	words := GetWordsFromLanguage(English)
	defer words.Free()

	assert.Equal(t, 2048, words.Length())
	assert.Equal(t, "abandon", words.Items[0])
	assert.Equal(t, "zoo", words.Items[2047])
}

func TestGenerateMnemonic(t *testing.T) {
	m, err := GenerateMnemonicRandom(English, Entropy128bits)
	defer m.Free()
	if err != nil {
		panic(err)
	}

	fmt.Println(m.Status)
}

func TestGenerateMnemonicFromBytes(t *testing.T) {
	entropy, err := hex.DecodeString("f0b9c942b9060af6a82d3ac340284d7e")
	assert.Nil(t, err)

	res, err := GenerateMnemonicFromBytes(entropy, English, Entropy128bits)
	defer res.Free()

	assert.Nil(t, err)
	assert.Equal(t, 12, len(res.Words))
	assert.Equal(t, "vague soft expose improve gaze kitten pass point select access battle wish", res.Raw)

	fmt.Println(res.Raw)

}

func TestGenerateMnemonicFromBytesInvalidEntropyLength(t *testing.T) {
	entropy, err := hex.DecodeString("f0b9c942b9060af6a82d3ac340284d7ef0b9c942b9060af6a82d3ac340284d7e")
	assert.Nil(t, err)

	res, err := GenerateMnemonicFromBytes(entropy, English, Entropy128bits)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), fmt.Sprintf("Invalid bytes size: must be equals to entropy - %d", Entropy128bits))
}

func TestWordsToSeed(t *testing.T) {
	words := "lock silly satisfy version solution bleak rain candy phone loan powder dose"
	expectedSeedHex := "ffd6a5e899b691e807718236da7f57932c76b33f402a652b8e216930ab3f72278e1d059e9b372db186662365e6e3bbdba689fbab5595e62fefa8e7687e11daa5"

	seed := MinterData64{}
	written := 0
	WordsToSeed(words, &seed, &written)

	hexed := seed.ToHexString()
	assert.Equal(t, expectedSeedHex, hexed)
	assert.Equal(t, len(seed.data), written)

	fmt.Println("Expected hex: ", expectedSeedHex)
	fmt.Println("Resulted hex: ", hexed)
}

func TestValidateMnemonicWords(t *testing.T) {
	words1 := "lock silly satisfy version solution bleak rain candy phone loan powder dose"
	res1 := ValidateMnemonicWords(English, words1)
	assert.True(t, res1)

	words2 := "wtf wtf wtf wtf wtf wtf wtf wtf wtf wtf wtf wtf"
	res2 := ValidateMnemonicWords(English, words2)
	assert.False(t, res2)

	words3 := "lock silly satisfy version solution wtf wtf wtf wtf wtf wtf wtf"
	res3 := ValidateMnemonicWords(English, words3)
	assert.False(t, res3)

	words4 := ""
	res4 := ValidateMnemonicWords(English, words4)
	assert.False(t, res4)
}

func TestMakeBip32RootKey(t *testing.T) {
	words := "lock silly satisfy version solution bleak rain candy phone loan powder dose"
	seed := MinterData64{}
	written := 0
	WordsToSeed(words, &seed, &written)

	rootKey := MakeBip32RootKey(seed)
	defer rootKey.Free()

	assert.Equal(t, uint32(0), rootKey.Index)
	assert.Equal(t, uint8(0x0), rootKey.Depth)
	assert.Equal(t, uint32(0), rootKey.Fingerprint)

	assert.Equal(t, 32, len(rootKey.PrivateKey.data))
	assert.NotEqual(t, uint8(0x0), rootKey.PrivateKey.data[0])
	assert.NotEqual(t, uint8(0x0), rootKey.PublicKey.data[0])
	assert.NotEqual(t, uint8(0x0), rootKey.PublicKey.data[32])
}

func TestMakeExtendedKey(t *testing.T) {
	words := "lock silly satisfy version solution bleak rain candy phone loan powder dose"
	seed := MinterData64{}
	written := 0
	WordsToSeed(words, &seed, &written)

	rootKey := MakeBip32RootKey(seed)
	//defer rootKey.Free()

	extKey := MakeExtendedKey(rootKey, "m/44'/60'/0'/0/0")
	//defer extKey.Free()

	expectedPrivateKey := "fd90261f5bd702ffbe7483c3b5aa7b76b1f40c1582cc6a598120b16067d3cb9a"

	assert.Equal(t, expectedPrivateKey, extKey.PrivateKey.ToHexString())

}
