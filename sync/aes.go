package sync

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/pkg/errors"
)

type AES struct {
	key       []byte
	nonce     []byte
	blockSize int
}

func NewAES(base64Key, base64Nonce string) (*AES, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, errors.Wrap(err, "decode key")
	}

	nonce, err := base64.StdEncoding.DecodeString(base64Nonce)
	if err != nil {
		return nil, errors.Wrap(err, "decode IV")
	}

	return &AES{
		key:       key,
		nonce:     nonce,
		blockSize: aes.BlockSize,
	}, nil
}

func (h *AES) Encrypt(input []byte) ([]byte, error) {
	block, err := aes.NewCipher(h.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, h.nonce, input, nil), nil
}

func (h *AES) Decrypt(input []byte) ([]byte, error) {
	block, err := aes.NewCipher(h.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, h.nonce, input, nil)
}

func (h *AES) Validate() (bool, error) {
	if len(h.key) != 32 {
		return false, errors.New("key is 32 bytes (256 bits)")
	}
	if len(h.nonce) != 12 {
		return false, errors.New("nonce is 12 bytes (96 bits)")
	}
	return true, nil
}
