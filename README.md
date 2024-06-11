# Keystone Signing Backend

## Description

This repository holds the code for a simple backend service to support air-gapped transaction signing using the a HD Wallet.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/crustyapples/go-keystone-backend.git
```

2. Change to the project directory:

```bash
cd go-keystone
```

3. Install dependencies:

```bash
go mod download
```

## Usage

1. Build the project: `go build`
2. Run the project: `go run main.go`

## API Docs

### Get Fingerprint

`POST /get-fingerprint`

### Request Body

```json
{
  "urData": "string"
}
```

### Response Body

```json
{
  "sourceFingerprint": "string"
}
```

### Get Sign Request

`POST /get-sign-request`

### Request Body

```json
{
  "txData": {
    "Nonce": "uint64",
    "To": "string",
    "Value": "uint64",
    "GasLimit": "uint64",
    "GasPrice": "uint64",
    "Data": "string",
    "ChainID": "int"
  },
  "fingerprint": "string"
}

```

### Response Body

```json
{
  "ethSignRequestCbor": "byte array"
}
```

### Sign Transaction

`POST /sign-transaction`

### Request Body

```json
{
  "signature": "string",
  "txData": {
    "Nonce": "uint64",
    "To": "string",
    "Value": "uint64",
    "GasLimit": "uint64",
    "GasPrice": "uint64",
    "Data": "string",
    "ChainID": "int"
  },
  "signer": "string"
}

```

### Response

```json
{
"signedTxn": "string"
}
```

## License

This project is licensed under the [MIT License](LICENSE).
