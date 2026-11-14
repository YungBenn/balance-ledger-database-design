package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

// PASETO v2 local tokens are encrypted with XChaCha20-Poly1305, which needs a 32-byte key.
const symmetricKeySize = 32

var symmetricKey []byte

// SetSymmetricKey sets the key used to create and verify tokens. Call it once at startup.
func SetSymmetricKey(key string) error {
	if len(key) != symmetricKeySize {
		return fmt.Errorf("TOKEN_SYMMETRIC_KEY must be exactly %d bytes, got %d", symmetricKeySize, len(key))
	}
	symmetricKey = []byte(key)
	return nil
}

func Create(userID string, email string, fullName string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, email, fullName, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := paseto.NewV2().Encrypt(symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

func Verify(token string) (*Payload, error) {
	payload := &Payload{}

	err := paseto.NewV2().Decrypt(token, symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
