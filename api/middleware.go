package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dawitfrazer/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	autorizationHeaderKey  = "authorization"
	autorizationTypeBearer = "bearer"
	autorizationPaylodKey  = "authorization_payload"
)

func authMiddleware(tokenMake token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationHeader := ctx.GetHeader(autorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid autorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		autorizationType := strings.ToLower(fields[0])
		if autorizationType != autorizationTypeBearer {
			err := fmt.Errorf("unsupported autorization type %s", autorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMake.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(autorizationPaylodKey, payload)
		ctx.Next()
	}
}
