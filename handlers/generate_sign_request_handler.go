package handlers

import (
	// "encoding/hex"

	"encoding/binary"
	"encoding/hex"
	"fmt"
	"go-keystone/mod/crypto"
	"go-keystone/mod/utils"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"

	// "github.com/fxamacker/cbor/v2"
	"github.com/gofiber/fiber/v2"
)

type GenerateSignRequestData struct {
	TxData      crypto.TxParams `json:"txData"`
	Fingerprint string          `json:"fingerprint"`
}

func GenerateSignRequestHandler(c *fiber.Ctx) error {
	var request GenerateSignRequestData

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// fmt.Println("request: ", request)

	txData := request.TxData
	fingerprint := request.Fingerprint

	if fingerprint == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fingerprint is required"})
	}

	// conver txData value from Int to big.Int
	value := new(big.Int).SetUint64(txData.Value)
	gasPrice := new(big.Int).SetUint64(txData.GasPrice)
	tx := types.NewTransaction(txData.Nonce, common.HexToAddress(txData.To), value, txData.GasLimit, gasPrice, common.FromHex(txData.Data))

	txRlp, err := rlp.EncodeToBytes(tx)
	// convert to hex
	txRlpHex := hex.EncodeToString(txRlp)
	fmt.Println("txRlpHex: ", txRlpHex)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// generate requestID using uuid v4
	requestID := uuid.New().String()
	// split the uuid by '-' and then concatenate to form a single string
	requestID = requestID[:8] + requestID[9:13] + requestID[14:18] + requestID[19:23] + requestID[24:]
	fmt.Println(requestID)

	// convert string to uint8array to byte array
	requestIDBytes, _ := hex.DecodeString(requestID)
	fmt.Println(requestIDBytes)

	origin := "metamask"
	var dataType crypto.DataType = crypto.Transaction

	fingerprintBytes, _ := hex.DecodeString(fingerprint)
	// convert fingerprintBytes to int
	fingerprintInt := binary.BigEndian.Uint32(fingerprintBytes)

	derivationPath := crypto.KeyPath{
		Components: []interface{}{
			// hardcoded derivation path, should ideally be taken in from the hdkey
			44, true, 60, true, 0, true, 0, false, 0, false,
		},
		SourceFingerprint: fingerprintInt,
	}

	ethSignRequest := crypto.EVMSignRequest{
		RequestID:      requestIDBytes,
		SignData:       txRlp,
		DataType:       dataType,
		ChainID:        txData.ChainID,
		DerivationPath: derivationPath,
		Origin:         origin,
	}

	// Create TagSet (safe for concurrency)
	tags := cbor.NewTagSet()

	// Register tags
	tags.Add(
		cbor.TagOptions{EncTag: cbor.EncTagRequired, DecTag: cbor.DecTagRequired},
		reflect.TypeOf(crypto.RequestIDType{}), 37)
	tags.Add(
		cbor.TagOptions{EncTag: cbor.EncTagRequired, DecTag: cbor.DecTagRequired},
		reflect.TypeOf(crypto.KeyPath{}), 304)

	ethSignRequestBytes, _ := utils.ToCBOR(&ethSignRequest, tags)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ethSignRequestCbor": ethSignRequestBytes})
}
