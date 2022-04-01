package sync

import (
	"encoding/base64"
	"testing"

	"gotest.tools/assert"
)

func TestAES_EncryptDecrypt(t *testing.T) {
	key := "12345678901234567890123456789012"
	nonce := "123456789012"
	plainText := []byte("test value")

	base64Key := base64.StdEncoding.EncodeToString([]byte(key))
	base64Nonce := base64.StdEncoding.EncodeToString([]byte(nonce))

	h, err := NewAES(base64Key, base64Nonce)
	assert.NilError(t, err)

	valid, err := h.Validate()
	assert.Equal(t, valid, true)
	assert.NilError(t, err)

	encrypted, err := h.Encrypt(plainText)
	assert.NilError(t, err)

	decrypted, err := h.Decrypt(encrypted)
	assert.NilError(t, err)
	assert.Equal(t, string(plainText), string(decrypted))
}
