package handlers

import (
	"encoding/hex"
	"fmt"

	"go-keystone/mod/crypto"

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

	typ, cbor, err := decoder.Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("typ: ", typ)
	fmt.Println("cbor: ", cbor)

	// ethSignature, err := eth.NewEthSignatureFromCBOR(cbor)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// }

	// sig := ethSignature.Signature()
	// r := sig[:32]
	// s := sig[32:64]
	// v := sig[64]

	// value := new(big.Int).SetUint64(txData.Value)
	// tx := types.NewTransaction(txData.Nonce, common.HexToAddress(txData.To), value, txData.GasLimit, txData.GasPrice, common.FromHex(txData.Data))
	// txWithSig, err := tx.WithSignature(types.HomesteadSigner{}, append(r, append(s, v)...))
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// }

	txWithSig := "txWithSig"

	signedTxn, err := rlp.EncodeToBytes(txWithSig)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"signedTxn": hex.EncodeToString(signedTxn)})
}
