// Package helpers
package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

const (
	rsaKeySize     = 256
	nonceSize      = 12
	staticFullSize = rsaKeySize + nonceSize
)

// LoadPublicKey - Load the server's public key from PEM
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	pubIfc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubIfc.(*rsa.PublicKey), nil
}

// EncryptForServer Encrypt a short message for the server
func EncryptForServer(pub *rsa.PublicKey, message []byte) ([]byte, error) {

	// - сообщения слишком маленькие, пришлось переписать
	//return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, message, nil)

	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("generate AES key: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("rand read nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, message, nil)

	encryptedAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("encrypt AES kew: %w", err)
	}

	final := make([]byte, 0, len(encryptedAesKey)+len(nonce)+len(ciphertext)) // capacity задаём, но длина = 0
	final = append(final, encryptedAesKey...)
	final = append(final, nonce...)
	final = append(final, ciphertext...)
	return final, err
}

// LoadPrivateKey On the server, load the private key and decrypt:
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func DecryptFromClient(priv *rsa.PrivateKey, encrypted []byte) ([]byte, error) {
	//return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
	if len(encrypted) < staticFullSize {
		return nil, errors.New("encrypted payload size")
	}
	encryptedAESKey := encrypted[:rsaKeySize]             // первые 256
	nonce := encrypted[rsaKeySize : rsaKeySize+nonceSize] // 12 после этого
	ciphertext := encrypted[rsaKeySize+nonceSize:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encryptedAESKey, nil)

	if err != nil {
		return nil, fmt.Errorf("decrypt OAEP: 5: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("new AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm open: %w", err)
	}

	return plaintext, nil
}
