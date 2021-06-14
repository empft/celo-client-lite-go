package celotoken

import "math/big"

type CeloToken string

const (
	// CELO - previously known as cGLD
	CELO CeloToken = "CELO"
	// CUSD - Celo Dollar
	CUSD CeloToken = "cUSD"
	// CEUR - Celo Euro
	CEUR CeloToken = "cEUR"
)

type BalanceMap map[CeloToken]*big.Int

type Balance struct {
	CELO *big.Int
	CUSD *big.Int
	CEUR *big.Int
}