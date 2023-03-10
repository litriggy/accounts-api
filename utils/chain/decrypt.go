package chain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mergermarket/go-pkcs7"
)

func Decrypt(encrypted string, secPw string) (string, error) {
	CIPHER_KEY := fmt.Sprintf("%x", md5.Sum([]byte(secPw)))

	fmt.Println(CIPHER_KEY)

	key := []byte(CIPHER_KEY)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	fmt.Printf("%x", iv)
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	//fmt.Printf(" ciphered: [ %s ]", cipherText)
	privateKey := string(cipherText)
	return privateKey, nil
}
