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

	// "github.com/fxamacker/cbor/v2"
	"github.com/gofiber/fiber/v2"
)

type GenerateQRRequest struct {
	TxData      crypto.TxParams `json:"txData"`
	Fingerprint string          `json:"fingerprint"`
}

func GenerateQRHandler(c *fiber.Ctx) error {
	var request GenerateQRRequest
	var qrData []string
	var signRequest crypto.EVMSignRequest

	testSignRequestHex := "a6015011fb0bb6192a4e6c93b1cbb0ba5826f30259026ef9026b04850f485c6b8e83061a8094bb60d7403f488a14fc33d926e7bec482291ae07c80b902446a7612020000000000000000000000003ebd5f2ae6877023b4330ab9a6aeb4a43ec0181f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000044a9059cbb000000000000000000000000b9fda735e6572c26cef89f68e95b0423d895098500000000000000000000000000000000000000000000000000000000000f4240000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041000000000000000000000000a44d8ffd5b3864648ccb32bdff04601b3cea783d000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008080800301041a001103e705a2018a182cf5183cf500f500f400f4021a5722f47e07686d6574616d61736b"

	raw, _ := hex.DecodeString(testSignRequestHex)
	utils.FromCBOR(raw, &signRequest, nil)

	// fmt.Println("fingerprint: ", *signRequest.DerivationPath.SourceFingerprint)

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

	requestID := "11fb0bb6192a4e6c93b1cbb0ba5826f3"
	// convert string to uint8array to byte array
	requestIDBytes, _ := hex.DecodeString(requestID)

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
	ethSignRequestHex := hex.EncodeToString(ethSignRequestBytes)
	fmt.Println("ethSignRequestHex: ", ethSignRequestHex)

	encoder := utils.NewFountainEncoder(ethSignRequestBytes, 200, 0, 200)

	// Generate the encoded parts
	for !encoder.IsComplete() {
		part := encoder.NextPart()
		if err != nil {
			fmt.Println("Error generating next part:", err)
		}
		partEncoded := utils.EncodePart("eth-sign-request", part)
		if err != nil {
			fmt.Println("Error encoding part:", err)
		}
		fmt.Println(partEncoded)
		qrData = append(qrData, partEncoded)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"qrData": qrData})
}
