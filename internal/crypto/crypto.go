package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var apiKey = []byte{
	0xEE, 0xBC, 0x1F, 0x57, 0x48, 0x7F, 0x51, 0x92,
	0x1C, 0x04, 0x65, 0x66, 0x5F, 0x8A, 0xE6, 0xD1,
	0x65, 0x8B, 0xB2, 0x6D, 0xE6, 0xF8, 0xA0, 0x69,
	0xA3, 0x52, 0x02, 0x93, 0xA5, 0x72, 0x07, 0x8F,
}

const nonceSize = 12

var ErrCiphertextTooShort = errors.New("crypto: ciphertext too short")

func newGCM() (cipher.AEAD, error) {
	block, err := aes.NewCipher(apiKey)
	if err != nil {
		return nil, fmt.Errorf("crypto: invalid AES key: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: GCM init failed: %w", err)
	}
	return gcm, nil
}

func Encrypt(plaintext []byte) (string, error) {
	gcm, err := newGCM()
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("crypto: read random nonce: %w", err)
	}

	sealed := gcm.Seal(nonce[:nonceSize:nonceSize], nonce, plaintext, nil)

	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func Decrypt(raw []byte) ([]byte, error) {
	if len(raw) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	gcm, err := newGCM()
	if err != nil {
		return nil, err
	}

	nonce := raw[:nonceSize]
	ciphertext := raw[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("crypto: decryption failed: %w", err)
	}
	return plaintext, nil
}
