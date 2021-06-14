package contractkit

import (
	"context"

	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/ethclient"
	"github.com/celo-org/celo-blockchain/rpc"
	"gitlab.com/stevealexrs/celo-client-lite-go/contractkit/celotoken"
	"gitlab.com/stevealexrs/celo-client-lite-go/contractkit/contracts"
	"gitlab.com/stevealexrs/celo-client-lite-go/contractkit/registry"
)

type Kit struct {
	Registry registry.Registry
	GoldToken *contracts.GoldToken
	StableToken *contracts.StableToken
	StableTokenEUR *contracts.StableToken
}

func New(url string) (*Kit, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)
	return NewWithEthClient(ethClient)
}

func NewWithEthClient(eth *ethclient.Client) (*Kit, error) {
	reg, err := registry.New(eth)
	if err != nil {
		return nil, err
	}

	goldToken, err := reg.GetGoldTokenContract(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	stableToken, err := reg.GetStableTokenContract(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	stableTokenEUR, err := reg.GetStableTokenEURContract(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &Kit{
		Registry: reg,
		GoldToken: goldToken,
		StableToken: stableToken,
		StableTokenEUR: stableTokenEUR,
	}, nil
}

// Get the balance of all tokens including native currency token
func (k *Kit) Balance(address string) (*celotoken.Balance, error) {
	addr := common.HexToAddress(address)

	gold, err := k.GoldToken.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		return nil, err
	}

	usd, err := k.StableToken.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		return nil, err
	}

	eur, err := k.StableTokenEUR.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		return nil, err
	}

	balance := &celotoken.Balance{
		CELO: gold,
		CUSD: usd,
		CEUR: eur,
	}
	return balance, nil
}