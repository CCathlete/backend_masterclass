package api

import (
	mockdb "backend-masterclass/db/mock"
	"backend-masterclass/token"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Adds an authorization header that contains <tokenType> <token> to the request.
func addAuthorisation(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {

	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	request.Header.Set(authorisationHeaderKey,
		fmt.Sprintf("%s %s", authorizationType, token))
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(
				t *testing.T,
				request *http.Request,
				tokenMaker token.Maker,
			) {
				addAuthorisation(
					t,
					request,
					tokenMaker,
					authorisationTypeBearer,
					"user",
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			// The client doesn't provide any authorization header.
			name: "No Authorization",
			setupAuth: func(
				t *testing.T,
				request *http.Request,
				tokenMaker token.Maker,
			) {
				// Do nothing.
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			// Here we provide a token type that is not supported.
			name: "Unsupported Authorization",
			setupAuth: func(
				t *testing.T,
				request *http.Request,
				tokenMaker token.Maker,
			) {
				unsupportedTokenType := "unsupported"

				addAuthorisation(
					t,
					request,
					tokenMaker,
					unsupportedTokenType,
					"user",
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			// Invalid format means that the token has less than 2 fields, i.e <tokenType> <token> where one of these is an empty string.
			name: "Invalid Authorization Format",
			setupAuth: func(
				t *testing.T,
				request *http.Request,
				tokenMaker token.Maker,
			) {
				invalidTokenType := ""

				addAuthorisation(
					t,
					request,
					tokenMaker,
					invalidTokenType,
					"user",
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Expired Token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			// ------------Setting up main engine resources-----------------
			// We don't need a db in this test group.
			var testStore *mockdb.MockStore = nil
			// Token maker is set up in this function.
			server := newTestServer(t, testStore)

			// ------------Setting up fake routing-----------------
			testHandlerFunc := func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			}

			authPath := "/auth"

			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				testHandlerFunc,
			)

			// ----------Setting up fake request & response writer-------------
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)

			// ------------------Running the test------------------------------
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
