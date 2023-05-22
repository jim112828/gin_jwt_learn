package main

import (
	"context"
	"log"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

var (
	signingKey = []byte("secret")

	issuer = "go-jwt-middleware-example"

	audience = []string{"audience-example"}

	keyFunc = func(ctx context.Context) (interface{}, error) {

		return signingKey, nil

	}

	customClaims = func() validator.CustomClaims {
		return &CustomClaimExample{}
	}
)

func checkJWT() gin.HandlerFunc {

	jwtValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		issuer,
		audience,
		validator.WithCustomClaims(customClaims),
		validator.WithAllowedClockSkew(30*time.Second),
	)

	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx *gin.Context) {
		encounteredError := true

		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			ctx.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid"},
			)
		}

	}

}
