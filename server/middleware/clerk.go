package middleware

import (
	"context"
	"log"
	"server/config"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
)

type ClerkUser struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"primary_email_address"`
}

func ClerkMiddleware() fiber.Handler {
	clerkSecretKey := config.Config("CLERK_SECRET_KEY")
	clerkPublishableKey := config.Config("CLERK_PUBLISHABLE_KEY")

	if clerkSecretKey == "" {
		panic("CLERK_SECRET_KEY environment variable is required")
	}
	if clerkPublishableKey == "" {
		panic("CLERK_PUBLISHABLE_KEY environment variable is required")
	}

	// Initialize Clerk client with secret key
	clerk.SetKey(clerkSecretKey)

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

		// Validate token with Clerk JWT verification using publishable key
		claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			log.Printf("Token verification failed for token ending in ...%s: %v",
				sessionToken[len(sessionToken)-4:], err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"debug": "Token verification failed - please refresh your session",
			})
		}

		// Extract user ID from claims
		userID := claims.Subject
		if userID == "" {
			log.Printf("Token verified but missing user ID in claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token: missing user ID",
			})
		}

		log.Printf("Successfully authenticated user: %s", userID)

		c.Locals("clerkUserID", userID)
		c.Locals("clerkSessionToken", sessionToken)
		c.Locals("clerkClaims", claims)

		return c.Next()
	}
}

func GetClerkUserIDFromContext(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals("clerkUserID").(string)
	if !ok || userID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}

	return userID, nil
}

func OptionalClerkMiddleware() fiber.Handler {
	clerkSecretKey := config.Config("CLERK_SECRET_KEY")
	clerkPublishableKey := config.Config("CLERK_PUBLISHABLE_KEY")

	if clerkSecretKey == "" || clerkPublishableKey == "" {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Initialize Clerk client with secret key
	clerk.SetKey(clerkSecretKey)

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
			claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
				Token: sessionToken,
			})
			if err == nil && claims.Subject != "" {
				c.Locals("clerkUserID", claims.Subject)
				c.Locals("clerkSessionToken", sessionToken)
				c.Locals("clerkClaims", claims)
			}
		}

		return c.Next()
	}
}

func GetClerkSecretKey() string {
	return config.Config("CLERK_SECRET_KEY")
}

func GetClerkPublishableKey() string {
	return config.Config("CLERK_PUBLISHABLE_KEY")
}
