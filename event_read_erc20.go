package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	token "hello/contracts_erc20"
	"log"
	"math/big"
	"strings"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

func main() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
	}

	contractAddr := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6383820),
		ToBlock:   big.NewInt(6383820),
		Addresses: []common.Address{
			contractAddr,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)

	for _, vLog := range logs {
		fmt.Printf("log: %s\n", vLog.BlockNumber)
		fmt.Printf("log: %s\n", vLog.Index)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("log: %s\n", vLog.BlockNumber)
			var transferEvent LogTransfer

			err := contractAbi.Unpack(&transferEvent, vLog)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf(transferEvent.From.Hex())
			fmt.Printf(transferEvent.To.Hex())
			fmt.Printf(transferEvent.Tokens.String())

		case logTransferSigHash.Hex():
			fmt.Printf("Log name: Approval\n")
			var approveEvent LogApproval

			_, err := contractAbi.Unpack(&approveEvent, "Approval", vLog)
			if err != nil {
				log.Fatal(err)
			}

			approveEvent.TokenOwner = common.HexToAddress(vLog.Topics[1].Hex())
			approveEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf(approveEvent.TokenOwner.Hex())
			fmt.Printf(approveEvent.Spender.Hex())
			fmt.Printf(approveEvent.Tokens.String())

		}
		fmt.Println("\n\n")
	}
}
