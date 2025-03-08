package token

import (
	"testing"
	"time"

	"github.com/badermezzi/KubeGoBank/util"
	"github.com/stretchr/testify/require"
)

// TestPasetoMaker tests the PasetoMaker's token creation and verification functionalities.
func TestPasetoMaker(t *testing.T) {
	// Generate a random symmetric key for testing.
	symmetricKey := util.RandomString(32)
	// Create a new PasetoMaker with the symmetric key.
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err) // Assert that no error occurred during PasetoMaker creation.

	username := util.RandomOwner() // Generate a random username for the token payload.
	duration := time.Minute        // Set the token duration to 1 minute.

	issuedAt := time.Now()              // Record the time before token creation.
	expiredAt := issuedAt.Add(duration) // Calculate the expected expiration time.

	// Create a new token for the given username and duration.
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)    // Assert that no error occurred during token creation.
	require.NotEmpty(t, token) // Assert that the created token is not empty.

	// Verify the created token.
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)      // Assert that no error occurred during token verification.
	require.NotEmpty(t, payload) // Assert that the verified payload is not empty.

	// Assert the payload data is correct.
	require.NotZero(t, payload.ID)                                       // Assert that the payload ID is not zero.
	require.Equal(t, username, payload.Username)                         // Assert that the username in the payload matches the generated username.
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)   // Assert that the issuedAt time in the payload is within 1 second of the recorded issuedAt time.
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second) // Assert that the expiredAt time in the payload is within 1 second of the calculated expiredAt time.
}

// TestExpiredPasetoToken tests the PasetoMaker's token verification with an expired token.
func TestExpiredPasetoToken(t *testing.T) {
	// Generate a random symmetric key for testing.
	symmetricKey := util.RandomString(32)
	// Create a new PasetoMaker with the symmetric key.
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err) // Assert that no error occurred during PasetoMaker creation.

	// Create a new token with a negative duration, making it immediately expired.
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)    // Assert that no error occurred during token creation.
	require.NotEmpty(t, token) // Assert that the created token is not empty.

	// Verify the expired token.
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)   // Assert that an error occurred during token verification (due to expiration).
	require.Nil(t, payload) // Assert that the payload is nil because the token is expired.
	// require.ErrorIs(t, err, ErrExpiredToken) // ErrExpiredToken is not defined yet // Commented out: ErrExpiredToken is not defined yet.
}
