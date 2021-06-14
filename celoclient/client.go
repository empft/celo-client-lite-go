package celoclient

import (
	"math/big"
	"time"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"gitlab.com/stevealexrs/celo-client-lite-go/connection"
	"gitlab.com/stevealexrs/celo-client-lite-go/contractkit"
	"gitlab.com/stevealexrs/celo-client-lite-go/contractkit/celotoken"
)

type Client interface {
	Balance(string) (*celotoken.Balance, error)
	SendRawTransaction(string) (common.Hash, error)
	WaitForTransaction(common.Hash, time.Duration) (*types.Receipt, error)
	LatestBlock() (*big.Int, error)
}

type client struct {
	Connection *connection.Client
	Kit		   *contractkit.Kit
}

func New(url string) (Client, error) {
	conn, err := connection.New(url)
	if err != nil {
		return nil, err
	}

	kit, err := contractkit.NewWithEthClient(conn.Eth)
	if err != nil {
		return nil, err
	}

	return &client{Connection: conn, Kit: kit}, nil
}

// Get the balance of all tokens including native currency token
func (c *client) Balance(address string) (*celotoken.Balance, error) {
	return c.Kit.Balance(address)
}

// Send hex encoded signed transaction payload and returns the hash
func (c *client) SendRawTransaction(rawTx string) (common.Hash, error) {
	return c.Connection.SendRawTransaction(rawTx)
}

// Wait for the transaction to complete by polling the blockchain every 500ms
func (c *client) WaitForTransaction(txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	return c.Connection.WaitForTransaction(txHash, timeout)
}
func (c *client) LatestBlock() (*big.Int, error) {
	return c.Connection.LatestBlock()
}

