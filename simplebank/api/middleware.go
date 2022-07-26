package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ebaudet/simplebank/token"
	"github.com/gin-gonic/gin"
)

// This file contains the code for the middleware used in the API.
// The middleware is used to authenticate the user and to extract the token from the request.
//                     ________________________
//    send request    |         ROUTE          |
// -----------------> |    /accounts/create    |
//                    |________________________|
//                                |
//                     ___________V____________
//    ctx.Abort()     |        MIDDLEWARES     |
// <------------------| Logger(ctx), Auth(ctx) |
//    send response   |________________________|
//                                |  ctx.Next()
//                     ___________V__________
//    send response   |        HANDLER       |
// <------------------|  createAccount(ctx)  |
//                    |______________________|

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		tokenString := fields[1]
		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// add payload to the context
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
