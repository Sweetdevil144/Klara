package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ClerkMiddleware() fiber.Handler {
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		panic("CLERK_SECRET_KEY environment variable is required")
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		sessionToken := tokenParts[1]
		if sessionToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Empty session token",
			})
		}

		c.Locals("clerkSessionToken", sessionToken)

		return c.Next()
	}
}

func GetClerkUserIDFromContext(c *fiber.Ctx) (string, error) {
	token, ok := c.Locals("clerkSessionToken").(string)
	if !ok || token == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	return token, nil
}

func OptionalClerkMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Next()
		}

		sessionToken := tokenParts[1]
		if sessionToken != "" {
			c.Locals("clerkSessionToken", sessionToken)
		}

		return c.Next()
	}
}

func GetClerkSecretKey() string {
	return os.Getenv("CLERK_SECRET_KEY")
}
