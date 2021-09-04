package middlewares

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte(DotEnvVariable("JWT_SECRET"))

// IsAuthorized -> verify jwt header
func IsAuthorized(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				AuthorizationResponse("Invalid JWT token", w)
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			}
		} else {
			AuthorizationResponse("Not Authorized", w)
		}
	})
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT -> generate jwt
func GenerateJWT(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
