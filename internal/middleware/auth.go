package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type AuthConfig struct {
	SecretKey []byte
}

func DefaultConfig(secretKey string) *AuthConfig {
	return &AuthConfig{
		SecretKey: []byte(secretKey),
	}
}

func AuthMiddleware(config *AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "authorization header is required"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token format. Use 'Bearer <token>'"})
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return config.SecretKey, nil
			})
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				sub, ok := claims["sub"]
				if !ok {
					return c.JSON(http.StatusUnauthorized, echo.Map{"error": "sub not found in token"})
				}

				c.Set("sub", sub)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}
	}
}
