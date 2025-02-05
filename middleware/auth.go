package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	jwtSec := os.Getenv("JWT_SECRET")
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(403).JSON(fiber.Map{"message":"No cookie"})
	}

	tokenString = tokenString[7:]
		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure that the token is signed with the correct method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSec), nil
		})
	
		fmt.Println(token.Claims)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}
	

	// Extract user_id from the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		role := claims["role"].(string)
		// Attach the user_id to the context
		c.Locals("user_id", userID)
		c.Locals("role", role)
	}
		// If the token is valid, continue with the request
		return c.Next()
}