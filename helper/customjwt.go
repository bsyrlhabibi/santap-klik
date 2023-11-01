package helper

import (
	"santapKlik/configs"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func CreateToken(userID int, username, role string, cfg *configs.ProgramConfig) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JWTMiddleware(cfg *configs.ProgramConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]string{"message": "Token tidak ada"})
			}

			// Check for the "Bearer " prefix
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return c.JSON(401, map[string]string{"message": "Token tidak valid"})
			}

			// Extract the token part without the prefix
			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.Secret), nil
			})

			if err != nil {
				return c.JSON(403, map[string]string{"message": "Token tidak valid"})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				role, ok := claims["role"].(string)
				if !ok {
					return c.JSON(403, map[string]string{"message": "Token tidak valid"})
				}
				if role == "user" {
					c.Set("username", claims["username"])
					c.Set("role", role)
					return next(c)
				} else {
					return c.JSON(403, map[string]string{"message": "Anda tidak diizinkan mengakses rute ini"})
				}
			} else {
				return c.JSON(403, map[string]string{"message": "Token tidak valid"})
			}
		}
	}
}

func AdminMiddleware(cfg *configs.ProgramConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]string{"message": "Token tidak ada"})
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return c.JSON(401, map[string]string{"message": "Token tidak valid"})
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.Secret), nil
			})

			if err != nil {
				return c.JSON(403, map[string]string{"message": "Token tidak valid"})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				role, ok := claims["role"].(string)
				if !ok {
					return c.JSON(403, map[string]string{"message": "Token tidak valid"})
				}
				if role == "admin" {
					c.Set("username", claims["username"])
					c.Set("role", role)
					return next(c)
				} else {
					return c.JSON(403, map[string]string{"message": "Anda tidak diizinkan mengakses rute ini"})
				}
			} else {
				return c.JSON(403, map[string]string{"message": "Token tidak valid"})
			}
		}
	}
}
