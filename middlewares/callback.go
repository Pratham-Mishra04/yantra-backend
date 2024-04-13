package middlewares

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/gofiber/fiber/v2"
)

// TODO Check This
var dytePublicKey *rsa.PublicKey // Initialize this variable with your actual public key

func VerifyDyteWebHook(c *fiber.Ctx) error {
	//TODO Log that webhook was received
	fmt.Println("Webhook received")

	var reqBody struct{}

	c.BodyParser(&reqBody)

	// Get the signature from the request headers
	signature := c.Get("dyte-signature")
	if signature == "" {
		return &helpers.AppError{Code: fiber.StatusBadRequest, Message: "Missing signature"}
	}

	// Verify the signature
	hashed := sha256.Sum256([]byte(fmt.Sprintf("%v", reqBody)))
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return &helpers.AppError{Code: fiber.StatusBadRequest, Message: "Failed to decode signature", LogMessage: err.Error(), Err: err}
	}
	err = rsa.VerifyPKCS1v15(dytePublicKey, crypto.SHA256, hashed[:], signatureBytes)
	if err != nil {
		return &helpers.AppError{Code: fiber.StatusUnauthorized, Message: "Signature verification failed", LogMessage: err.Error(), Err: err}
	}

	return c.Next()
}
