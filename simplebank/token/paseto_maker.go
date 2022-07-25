package token

import (
	"fmt"
	"time"

	"github.com/vk-rv/pvx"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *pvx.ProtoV4Local
	symmetricKey *pvx.SymKey
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       pvx.NewPV4Local(),
		symmetricKey: pvx.NewSymmetricKey([]byte(symmetricKey), pvx.Version4),
	}
	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload)
	return token, payload, err
}

// VerifyToken checks if a token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey).ScanClaims(payload)
	if err != nil {
		if err == ErrExpiredToken {
			return nil, err
		}
		return nil, ErrInvalidToken
	}

	return payload, nil
}
