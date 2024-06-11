package handlers

import (
	"encoding/hex"
	"math/big"

	"go-keystone/mod/crypto"
	"go-keystone/mod/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gofiber/fiber/v2"
	"seedhammer.com/bc/ur"
)

type SignTransactionRequest struct {
	Signature string          `json:"signature"`
	TxData    crypto.TxParams `json:"txData"`
	Signer    string          `json:"signer"`
}

func SignTransactionHandler(c *fiber.Ctx) error {
	var request SignTransactionRequest
	var decoder ur.Decoder
	// var data crypto.TxParams

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	signature := request.Signature
	txData := request.TxData

	if signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "signature is required"})
	}

	if txData == (crypto.TxParams{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "txData is required"})
	}

	decoder.Add(signature)
	_, cborBytes, err := decoder.Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// decoding cbor to get the ETHSignature object
	var ethSignature crypto.ETHSignature
	err = utils.FromCBOR(cborBytes, &ethSignature, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	parsedSignature := ethSignature.Signature

	value := new(big.Int).SetUint64(txData.Value)
	gasPrice := new(big.Int).SetUint64(txData.GasPrice)
	tx := types.NewTransaction(txData.Nonce, common.HexToAddress(txData.To), value, txData.GasLimit, gasPrice, common.FromHex(txData.Data))
	txWithSig, err := tx.WithSignature(types.HomesteadSigner{}, parsedSignature[:65])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	signedTxn, err := rlp.EncodeToBytes(txWithSig)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"signedTxn": "0x" + hex.EncodeToString(signedTxn)})
}
