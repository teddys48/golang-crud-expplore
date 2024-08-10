package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func Decrypt(base64Ciphertext string, key []byte) ([]byte, error) {
	// Decode the base64 encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		return nil, err
	}

	// Create the AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := []byte("w%uC2u-80{t0M!h3")
	// ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	// Create the CBC mode decrypter
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the data
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	plaintext, err := pkcs7Unpad(ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func Encrypt(plaintext []byte, key []byte) (string, error) {
	// Ensure the key length is correct for AES-256 (32 bytes)
	if len(key) != 32 {
		return "", errors.New("key length must be 32 bytes for AES-256")
	}

	// Generate a random IV (16 bytes for AES)
	iv := []byte("w%uC2u-80{t0M!h3")

	// Create the AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Add PKCS#7 padding to the plaintext
	paddedPlaintext := pkcs7Pad(plaintext, aes.BlockSize)

	// Create CBC mode encrypter
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt the data
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Prepend IV to the ciphertext for use during decryption
	// ivAndCiphertext := append(iv, ciphertext...)

	// Encode the result as a base64 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func pkcs7Unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}

	padding := src[length-1]
	padLen := int(padding)

	if padLen == 0 || padLen > length {
		return nil, errors.New("invalid padding size")
	}

	for i := 0; i < padLen; i++ {
		if src[length-1-i] != padding {
			return nil, errors.New("invalid padding character")
		}
	}

	return src[:length-padLen], nil
}

func pkcs7Pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}
