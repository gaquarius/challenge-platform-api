package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var JWT_SECRET = []byte(DotEnvVariable("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// IsAuthorized -> verify jwt header
func IsAuthorized(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		//fmt.Println(authHeader)

		if len(authHeader) != 2 {
			AuthorizationResponse("Malformed JWT token", w)
		} else {
			jwtToken := authHeader[1]
			token, _ := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return JWT_SECRET, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "props", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				AuthorizationResponse("Unauthorized", w)
			}
		}
	})
}

// GenerateJWT -> generate jwt
func GenerateJWT(username string) (string, error) {
	//fmt.Println(`this is user`, username)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWT_SECRET)
	//fmt.Println(`this is it`, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
