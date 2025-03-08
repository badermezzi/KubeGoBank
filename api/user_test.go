package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/badermezzi/KubeGoBank/db/mock"
	db "github.com/badermezzi/KubeGoBank/db/sqlc"
	"github.com/badermezzi/KubeGoBank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore) db.User
		checkResponse func(recorder *httptest.ResponseRecorder, user db.User)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  "testuser",
				"password":  "password", // Plain password in test body
				"full_name": "Test User",
				"email":     "test@example.com",
			},
			buildStubs: func(store *mockdb.MockStore) db.User {
				hashedPassword, err := util.HashPassword("password") // Hash password for mock expectation
				require.NoError(t, err)

				arg := db.CreateUserParams{
					Username:       "testuser",
					HashedPassword: hashedPassword, // Use hashed password in mock arg
					FullName:       "Test User",
					Email:          "test@example.com",
				}

				user := db.User{
					Username:          arg.Username,
					HashedPassword:    arg.HashedPassword,
					FullName:          arg.FullName,
					Email:             arg.Email,
					PasswordChangedAt: time.Now(),
					CreatedAt:         time.Now(),
				}

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()). // Expect CreateUser call with hashed password
					Times(1).
					Return(user, nil)
				return user // Return user
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, user db.User) { // Accept user
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user) // Use user for comparison
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  "testuser",
				"password":  "password",
				"full_name": "Test User",
				"email":     "test@example.com",
			},
			buildStubs: func(store *mockdb.MockStore) db.User {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
				return db.User{} // Return empty user
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, user db.User) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidInput",
			body: gin.H{
				"username":  "", // invalid username
				"password":  "password",
				"full_name": "Test User",
				"email":     "test@example.com",
			},
			buildStubs: func(store *mockdb.MockStore) db.User {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0) // Expect no calls to CreateUser
				return db.User{} // Return empty user
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, user db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  "testuser",
				"password":  "password",
				"full_name": "Test User",
				"email":     "test@example.com",
			},
			buildStubs: func(store *mockdb.MockStore) db.User {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"}) // Unique violation error code
				return db.User{} // Return empty user
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, user db.User) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  "testuser",
				"password":  "password",
				"full_name": "Test User",
				"email":     "invalid-email", // Invalid email format
			},
			buildStubs: func(store *mockdb.MockStore) db.User {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0) // Expect no calls to CreateUser
				return db.User{} // Return empty user
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, user db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"username":  "testuser",
				"password":  "short", // Invalid password (too short)
				"full_name": "Test User",
				"email":     "test@example.com",
			},
			buildStubs: func(store *mockdb.MockStore) db.User {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0) // Expect no calls to CreateUser
				return db.User{} // Return empty user
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, user db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalid username", // Invalid username (contains space)
				"password":  "password",
				"full_name": "Test User",
				"email":     "test@example.com",
			},
			buildStubs: func(store *mockdb.MockStore) db.User {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0) // Expect no calls to CreateUser
				return db.User{} // Return empty user
			},
			checkResponse: func(recorder *httptest.ResponseRecorder, user db.User) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			user := tc.buildStubs(store) // Get user from buildStubs
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server := newTestServer(store)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder, user) // Pass user to checkResponse
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse createUserResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotResponse.Username)
	require.Equal(t, user.FullName, gotResponse.FullName)
	require.Equal(t, user.Email, gotResponse.Email)
	require.NotEmpty(t, gotResponse.CreatedAt)
	require.NotEmpty(t, gotResponse.PasswordChangedAt) // PasswordChangedAt should be default value, not empty if not set explicitly in createUserResponse
	require.WithinDuration(t, user.CreatedAt, gotResponse.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, gotResponse.PasswordChangedAt, time.Second)

	fmt.Println("response body:", body.String()) // for debugging
}
