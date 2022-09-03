package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dawitfrazer/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)
func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	autorizationType string,
	username string,
	duration time.Duration,
){
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", autorizationType, token)
	request.Header.Set(autorizationHeaderKey, authorizationHeader)
}
func TestAuthMiddleWare(t *testing.T) {
		testCases := []struct{
			name string
			setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
			checkResponse func (t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
			  name: "OK",
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
						addAuthorization(t, request,tokenMaker, autorizationTypeBearer, "user", time.Minute)
				},
				checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder){
						require.Equal(t, http.StatusOK, recorder.Code)
				},
		},
		{
			name: "NoAuthoization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
				
			},
			checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder){
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
	},
	{
		name: "UnSpportedAuthoization",
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
			addAuthorization(t, request,tokenMaker, "unsuppoted", "user", time.Minute)
		},
		checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
},
{
	name: "InvalidAuthoizationFormat",
	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
		addAuthorization(t, request,tokenMaker, "", "user", time.Minute)
	},
	checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder){
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
	},
},
{
	name: "ExpiredToken",
	setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker){
			addAuthorization(t, request,tokenMaker, autorizationTypeBearer, "user", -time.Minute)
	},
	checkResponse: func (t *testing.T, recorder *httptest.ResponseRecorder){
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
	},
},
	}
		for i :=range testCases {
			tc := testCases[i]

			t.Run(tc.name, func(t *testing.T){
					server := newTestServer(t, nil)

					authPath := "/auth"
					server.router.GET(
						authPath,
						authMiddleware(server.tokenMaker),
						func(ctx *gin.Context) {
							ctx.JSON(http.StatusOK, gin.H{})
						},
					)

					recorder := httptest.NewRecorder()
					request, err := http.NewRequest(http.MethodGet, authPath, nil)
					require.NoError(t, err)

					tc.setupAuth(t, request, server.tokenMaker)
					server.router.ServeHTTP(recorder, request)
					tc.checkResponse(t, recorder)
			})
		}


}