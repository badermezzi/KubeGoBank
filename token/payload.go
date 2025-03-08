package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Payload contains the payload data of the token.
// It includes the token ID, username, issued at time, and expired at time.
type Payload struct {
	ID        uuid.UUID `json:"id"`        // ID is the unique identifier of the token.
	Username  string    `json:"username"`  // Username is the username of the token owner.
	IssuedAt  time.Time `json:"issued_at"` // IssuedAt is the time when the token was issued.
	ExpiredAt time.Time `json:"expired_at"`// ExpiredAt is the time when the token will expire.
}

// NewPayload creates a new token payload with the given username and duration.
// It returns a Payload and an error.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom() // Generate a new random UUID for the token ID.
	if err != nil {
		return nil, err // Return error if UUID generation fails.
	}

	payload := &Payload{
		ID:        tokenID,               // Assign the generated token ID to the payload.
		Username:  username,              // Assign the given username to the payload.
		IssuedAt:  time.Now(),            // Set the issued at time to the current time.
		ExpiredAt: time.Now().Add(duration), // Set the expired at time to the current time plus the given duration.
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not.
// It returns an error if the token is expired or invalid.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) { // Check if the current time is after the token's expiration time.
		return errors.New("token has expired") // Return an error if the token has expired.
	}
	return nil // Return nil if the token is still valid.
}
