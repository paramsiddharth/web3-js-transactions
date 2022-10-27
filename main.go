package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	rand.Seed(time.Now().UnixNano())

	SELF := os.Getenv("WEB3_SELF_ADDRESS")
	PEER := os.Getenv("WEB3_PEER_ADDRESS")
	KEY := os.Getenv("WEB3_SELF_KEY")
	THEFT := os.Getenv("WEB3_PEER_KEY")

	if len(SELF) < 1 ||
		len(PEER) < 1 ||
		len(KEY) < 1 ||
		len(THEFT) < 1 {
		log.Fatalln("Error: Environment variables not properly configured.")
	}

	getUrl := func() string {
		envVar := os.Getenv("WEB3_PROVIDER")
		if envVar != "" {
			return envVar
		}
		return "ws://localhost:8545"
	}

	ec, err := ethclient.Dial(getUrl())
	if err != nil {
		panic(err)
	}

	web3, err := web3.NewWeb3(getUrl())
	if err != nil {
		panic(err)
	}

	count := 10
	for i := 1; i <= count; i++ {
		fmt.Printf("\n--- Iteration %d ---\n\n", i)
		makeTransaction(ec, web3, SELF, PEER, KEY, THEFT)
	}
}

func makeTransaction(
	ec *ethclient.Client,
	web3 *web3.Web3,
	SELF string,
	PEER string,
	KEY string,
	THEFT string,
) {
	fmt.Println("--- Before transaction ---")

	selfBal := getBal(web3, SELF)
	fmt.Printf("Self: %v ETH\n", selfBal)

	peerBal := getBal(web3, PEER)
	fmt.Printf("Peer: %v ETH\n", peerBal)

	fmt.Println("--- Transaction ---")

	spend := selfBal.Cmp(peerBal) >= 0
	amount := new(big.Float).Mul(big.NewFloat(rand.Float64()), If(spend, selfBal, peerBal))
	amt, _ := amount.Float64()
	weis := toWei(web3, amt)
	to := If(spend, PEER, SELF)
	// from := If(spend, SELF, PEER)
	key := If(spend, KEY, THEFT)
	gasPrice, err := web3.Eth.GasPrice()
	if err != nil {
		panic(err)
	}

	privKey, err := crypto.HexToECDSA(key[2:])
	if err != nil {
		panic(err)
	}

	pubKey := privKey.PublicKey

	fromAdd := crypto.PubkeyToAddress(pubKey)
	nonce, err := ec.PendingNonceAt(context.Background(), fromAdd)
	if err != nil {
		panic(err)
	}

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(to),
		weis,
		uint64(21000),
		big.NewInt(int64(gasPrice)),
		[]byte{},
	)

	signed, err := types.SignTx(tx, types.HomesteadSigner{}, privKey)
	if err != nil {
		panic(err)
	}

	err = ec.SendTransaction(context.Background(), signed)
	if err != nil {
		panic(err)
	}

	tn, _, err := ec.TransactionByHash(context.Background(), signed.Hash())
	if err != nil {
		panic(err)
	}

	rct, err := ec.TransactionReceipt(context.Background(), tn.Hash())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Amount: %v ETH\n", new(big.Float).Mul(amount, new(big.Float).SetInt64(func() int64 {
		if spend {
			return -1
		}
		return 1
	}())))
	fmt.Printf("Transaction: %v\n", signed.Hash())
	fmt.Printf("Block: %v\n", rct.BlockNumber)

	fmt.Println("--- After transaction ---")

	fmt.Printf("Self: %v ETH\n", getBal(web3, SELF))
	fmt.Printf("Peer: %v ETH\n", getBal(web3, PEER))
}

func getBal(web3 *web3.Web3, add string) *big.Float {
	bal, err := web3.Eth.GetBalance(common.HexToAddress(add), nil)
	if err != nil {
		panic(err)
	}

	return web3.Utils.FromWei(bal)
}

func toWei(web3 *web3.Web3, val float64) *big.Int {
	return web3.Utils.ToWei(val)
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
