package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

// PasetoMaker is a PASETO token maker.
type PasetoMaker struct {
	paseto       *paseto.V2 // paseto is the PASETO v2 instance used to create and verify tokens.
	symmetricKey []byte    // symmetricKey is the secret key used to sign and verify PASETO tokens.
}

// NewPasetoMaker creates a new PasetoMaker.
// It takes a symmetricKey string as input and returns a Maker interface and an error.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// Check if the symmetric key length is valid (must be 32 bytes for PASETO v2).
	if len(symmetricKey) != 32 {
		return nil, fmt.Errorf("invalid key size: must be exactly 32 characters")
	}

	// Create a new PasetoMaker instance.
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(), // Initialize the PASETO v2 instance.
		symmetricKey: []byte(symmetricKey), // Convert the symmetricKey string to a byte slice.
	}
	return maker, nil
}

// CreateToken creates a new PASETO token for the given username and duration.
// It takes a username and duration as input and returns a token string and an error.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// Create a new payload for the token.
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err // Return error if payload creation fails.
	}

	// Encrypt the payload using the symmetric key and return the token string.
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// VerifyToken verifies if the token is valid or not.
// It takes a token string as input and returns a Payload and an error.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	// Create an empty payload to store the decrypted data.
	payload := &Payload{}

	// Decrypt the token using the symmetric key and unmarshal the payload into the payload struct.
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err) // Return error if token decryption fails.
	}

	// Check if the token payload is valid (not expired).
	err = payload.Valid()
	if err != nil {
		return nil, err // Return error if token payload is invalid.
	}

	return payload, nil // Return the valid payload and no error.
}
