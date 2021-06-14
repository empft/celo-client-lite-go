package account

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/celo-org/celo-blockchain/crypto"
)

const (
	publicKeyPrefix = "04"
)

func GeneratePrivateKey() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	return hexutil.Encode(privateKeyBytes)[2:]
}

func NewPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(key)
}

func NewPublicKey(key string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyPrefix + key)
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPubkey(publicKeyBytes)
}

func MakeAddress(s string) (common.Address) {
	return common.HexToAddress(s)
}

func MustDerivePublicKey(privateKey string) string {
	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		panic(err)
	}
	return publicKey
}

func DerivePublicKey(privateKey string) (string, error) {
	privKey, err := NewPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	return hexutil.Encode(publicKeyBytes)[4:], nil
}

func MustDeriveAddress(publicKey string) string {
	address, err := DeriveAddress(publicKey)
	if err != nil {
		panic(err)
	}
	return address
}

func DeriveAddress(publicKey string) (string, error) {
	publicKeyECDSA, err := NewPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	
	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()[2:], nil
}

func MustCompressPublicKey(publicKey string) string {
	compressedKey, err := CompressPublickey(publicKey)
	if err != nil {
		panic(err)
	}
	return compressedKey
}

func CompressPublickey(publicKey string) (string, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyPrefix + publicKey)
	if err != nil {
		return "", err
	}

	publicKeyECDSA, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(crypto.CompressPubkey(publicKeyECDSA))[2:], nil
}

func MustDecompressPublicKey(compressedKey string) string {
	publicKey, err := DecompressPublicKey(compressedKey)
	if err != nil {
		panic(err)
	}
	return publicKey
}

func DecompressPublicKey(compressedKey string) (string, error) {
	compressedKeyBytes, err := hex.DecodeString(compressedKey)
	if err != nil {
		return "", err
	}

	publicKeyECDSA, err := crypto.DecompressPubkey(compressedKeyBytes)
	if err != nil {
		return "", err
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	return hexutil.Encode(publicKeyBytes)[4:], nil
}