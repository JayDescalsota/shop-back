package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type Manager struct {
	key []byte
}

func NewManager(key string) *Manager {
	k := []byte(key)
	if len(k) != 32 {
		panic("encryption key must be 32 bytes (256 bits)")
	}
	return &Manager{key: k}
}

func (m *Manager) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(m.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func (m *Manager) Decrypt(encoded string) ([]byte, error) {
	ciphertext, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(m.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
