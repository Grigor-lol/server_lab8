package jwt_auth

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

var Users = map[string]string{
	"admin": "admin",
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Admin    bool   `json:"admin"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var JwtKey = []byte("my_secret_key")

func VerifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := c.Value
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if request.Method == "DELETE" {
			if !claims.Admin {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		endpointHandler(writer, request)
	}
}
