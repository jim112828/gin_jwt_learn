package main

import (
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", checkJWT(), func(ctx *gin.Context) {
		claims, ok := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

		if !ok {
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				map[string]string{"message": "Failed to cast custom JWT claims to specific type"},
			)
			return
		}
		customerClaims, ok := claims.CustomClaims.(*CustomClaimExample)
		if !ok {
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				map[string]string{"message": "Username is JWT claims was empty"},
			)
			return
		}

		if len(customerClaims.Username) == 0 {
			ctx.AbortWithStatusJSON(
				http.StatusBadRequest,
				map[string]string{"message": "Username in JWT claims was empty"},
			)
			return
		}

		ctx.JSON(http.StatusOK, claims)
	})

	log.Print("Server listening on http://localhost:3000")
	if err := http.ListenAndServe("0.0.0.0:3000", router); err != nil {
		log.Fatalf("Tere was an error with the http server: %v", err)
	}

}
