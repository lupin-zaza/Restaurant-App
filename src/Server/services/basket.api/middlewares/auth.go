package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jurabek/basket.api/models"
	"github.com/jurabek/basket.api/oidc"
	"github.com/patrickmn/go-cache"
)

// CreateAuth creates new instance of Auth
func CreateAuth() *Auth {

	claimsToValidate := map[string]interface{}{}
	claimsToValidate["aud"] = "menu-api"
	claimsToValidate["iss"] = "http://demo.restaurant-identity"

	httpClient := oidc.JWKHttpClient{}

	verifier := oidc.JwtVerifier{
		Cache:            cache.New(5*time.Minute, 10*time.Minute),
		Authority:        "http://localhost:5000",
		ClaimsToValidate: claimsToValidate,
		HTTPClient:       &httpClient,
	}

	verifier.New()
	auth := Auth{
		JwtVerifier: &verifier,
	}

	return &auth
}

// Auth represents AuthMiddleware
type Auth struct {
	JwtVerifier *oidc.JwtVerifier
}

// AuthMiddleware provides for securing Handlers
func (auth *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.URL.Path, "/api/v1/items") {
			c.Next()
			return
		}
		authorizationHeader := c.Request.Header.Get("Authorization") //Grab the token from the header
		if authorizationHeader == "" {                               //Token is missing, returns with error code 403 Unauthorized
			auth.abortMiddleware(c, fmt.Errorf("token is missing"))
			return
		}
		// The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			auth.abortMiddleware(c, fmt.Errorf("bearer token invalid format"))
			return
		}
		ok, err := auth.JwtVerifier.ValidateToken(bearerToken[1])
		if !ok && err != nil {
			auth.abortMiddleware(c, err)
			return
		}
		c.Next()
	}
}

func (auth *Auth) abortMiddleware(c *gin.Context, err error) {
	httpError := models.NewHTTPError(http.StatusUnauthorized, err)
	c.AbortWithStatusJSON(http.StatusUnauthorized, httpError)
}
