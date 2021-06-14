// Manages connection to celo blockchain
package connection

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/ethclient"
	"github.com/celo-org/celo-blockchain/rpc"
	"gitlab.com/stevealexrs/celo-client-lite-go/account"
)

const (
	minGasLimit = uint64(21000)
)

type Client struct {
	rpc *rpc.Client
	Eth *ethclient.Client
}

func New(url string) (*Client, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}

	ethClient := ethclient.NewClient(rpcClient)

	return &Client{
		rpc: rpcClient,
		Eth: ethClient,
	}, nil
}

func (c *Client) LatestBlock() (*big.Int, error) {
	header, err := c.Eth.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return header.Number, nil
}


func (c *Client) Balance(address string) (*big.Int, error) {
	return c.Eth.BalanceAt(context.Background(), common.HexToAddress(address), nil)
}

func (c *Client) SendNative(privateKey string, from string, to string, amount *big.Int, gatewayFeeRecipient *string, gatewayFee *big.Int) (common.Hash, error) {
	nonce, err := c.Eth.PendingNonceAt(context.Background(), common.HexToAddress(from))
	if err != nil {
		return common.Hash{}, err
	}

	gasPrice, err := c.Eth.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	var feeRecipient *common.Address
	var address common.Address
	if gatewayFeeRecipient != nil {
		address = common.HexToAddress(*gatewayFeeRecipient)
		feeRecipient = &address
	} else {
		feeRecipient = nil
	}

	tx := types.NewTransaction(
		nonce, 
		common.HexToAddress(to),
		amount,
		minGasLimit,
		gasPrice,
		nil,
		feeRecipient,
		gatewayFee,
		nil,
	)

	chainID, err := c.Eth.NetworkID(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	key, err := account.NewPrivateKey(privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	if err != nil {
		return common.Hash{}, err
	}
	
	return c.SendTransaction(signedTx)
}

func (c *Client) SendTransaction(tx *types.Transaction) (common.Hash, error) {
	return tx.Hash(), c.Eth.SendTransaction(context.Background(), tx)
}

func (c *Client) SendRawTransaction(rawTx string) (common.Hash, error) {
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := c.Eth.SendRawTransaction(context.Background(), rawTxBytes)
	if err != nil {
		return common.Hash{}, err
	}
	return *hash, nil
}

func (c *Client) WaitForTransaction(txHash common.Hash, timeout time.Duration) (*types.Receipt, error) {
	step := time.Millisecond * 500
	start := time.Now()
	
	for {
		if time.Since(start) > timeout {
			break
		}

		receipt, err := c.Eth.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return nil, err
		}

		if receipt != nil {
			if receipt.Status == 1 {
				return receipt, nil
			}
			return receipt, errors.New("transaction failed to execute")
		}
		time.Sleep(step)
	}
	return nil, errors.New("transaction not found within timeout period")
}