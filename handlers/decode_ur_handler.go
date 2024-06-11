package handlers

import (
	"fmt"
	"go-keystone/mod/crypto"
	"go-keystone/mod/utils"

	"github.com/gofiber/fiber/v2"
	"seedhammer.com/bc/ur"
)

type DecodeURRequest struct {
	URData string `json:"urData"`
}

func DecodeURHandler(c *fiber.Ctx) error {
	var request DecodeURRequest
	var KeystoneHDKey crypto.HDKey

	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if URData is provided
	if request.URData == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "urData is required"})
	}

	// Decode the URData
	var decoder ur.Decoder
	decoder.Add(request.URData)

	_, cborBytes, err := decoder.Result()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert CBOR bytes to HDKey object
	err = utils.FromCBOR(cborBytes, &KeystoneHDKey, nil)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert SourceFingerprint to hex string
	sourceFingerprint := fmt.Sprintf("%x", KeystoneHDKey.Origin.SourceFingerprint)

	// Return the result
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"sourceFingerprint": sourceFingerprint})
}
