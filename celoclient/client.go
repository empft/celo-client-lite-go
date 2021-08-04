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

type Client struct {
	Connection *connection.Client
	Kit        *contractkit.Kit
}

func New(url string) (*Client, error) {
	conn, err := connection.New(url)
	if err != nil {
		return nil, err
	}

	kit, err := contractkit.NewWithEthClient(conn.Eth)
	if err != nil {
		return nil, err
	}

	return &Client{Connection: conn, Kit: kit}, nil
}

// Get the balance of all tokens including native currency token
func (c *Client) Balance(address string) (*celotoken.Balance, error) {
	return c.Kit.Balance(address)
}

// Send hex encoded signed transaction payload and returns the hash
func (c *Client) SendRawTransaction(rawTx string) (common.Hash, error) {
	return c.Connection.SendRawTransaction(rawTx)
}

// Wait for the transaction to complete by polling the blockchain every 500ms
func (c *Client) WaitForTransaction(txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	return c.Connection.WaitForTransaction(txHash, timeout)
}
func (c *Client) LatestBlock() (*big.Int, error) {
	return c.Connection.LatestBlock()
}
