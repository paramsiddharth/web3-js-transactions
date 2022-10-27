package main

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"

	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

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

	web3, err := web3.NewWeb3(func() string {
		envVar := os.Getenv("WEB3_PROVIDER")
		if envVar != "" {
			return envVar
		}
		return "ws://localhost:8545"
	}())
	if err != nil {
		log.Fatalln(err)
	}

	count := 10
	for i := 1; i <= count; i++ {
		makeTransaction(web3, SELF, PEER, KEY, THEFT)
	}
}

func makeTransaction(
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
	to := If(spend, PEER, SELF)
	// from := If(spend, SELF, PEER)
	key := If(spend, KEY, THEFT)
	gasPrice, err := web3.Eth.GasPrice()
	if err != nil {
		panic(err)
	}

	web3.Eth.SetAccount(key)
	tx, err := web3.Eth.SendRawTransaction(
		common.HexToAddress(to),
		toWei(web3, amt),
		21000,
		big.NewInt(int64(gasPrice)),
		nil,
	)
	if err != nil {
		panic(err)
	}

	// web3.Eth.NewEIP1559Tx()

	fmt.Printf("Transaction: %v\n", tx)

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
