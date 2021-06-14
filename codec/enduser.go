package codec

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/celo-org/celo-blockchain/params"
)

// Copied from https://github.com/ethereum/go-ethereum/issues/21221

func CeloToWei(celo *big.Float) *big.Int {
	truncInt, _ := celo.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", celo), ".")[1]
	fracStr += strings.Repeat("0", 18 - len(fracStr))
	fracInt, _ :=  new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func WeiToCelo(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)  //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236)  //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func StringToCelo(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236)  //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}