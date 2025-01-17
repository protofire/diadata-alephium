package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/diadata-org/diadata/pkg/dia/scraper/blockchain-scrapers/blockchains/ethereum/diaCoingeckoOracleService"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/diadata-org/diadata/pkg/dia"
	models "github.com/diadata-org/diadata/pkg/model"
	"github.com/diadata-org/diadata/pkg/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	key := utils.Getenv("PRIVATE_KEY", "")
	key_password := utils.Getenv("PRIVATE_KEY_PASSWORD", "")
	deployedContract := utils.Getenv("DEPLOYED_CONTRACT", "")
	blockchainNode := utils.Getenv("BLOCKCHAIN_NODE", "")
	sleepSeconds, err := strconv.Atoi(utils.Getenv("SLEEP_SECONDS", "120"))
	if err != nil {
		log.Fatalf("Failed to parse sleepSeconds: %v")
	}
	frequencySeconds, err := strconv.Atoi(utils.Getenv("FREQUENCY_SECONDS", "120"))
	if err != nil {
		log.Fatalf("Failed to parse frequencySeconds: %v")
	}
	chainId, err := strconv.ParseInt(utils.Getenv("CHAIN_ID", "1"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse chainId: %v")
	}
	numCoins, err := strconv.Atoi(utils.Getenv("NUM_COINS", "50"))
	if err != nil {
		log.Fatalf("Failed to parse numCoins: %v")
	}

	/*
	 * Setup connection to contract, deploy if necessary
	 */

	conn, err := ethclient.Dial(blockchainNode)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	auth, err := bind.NewTransactorWithChainID(strings.NewReader(key), key_password, big.NewInt(chainId))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	var contract *diaCoingeckoOracleService.DIACoingeckoOracle
	err = deployOrBindContract(deployedContract, conn, auth, &contract)
	if err != nil {
		log.Fatalf("Failed to Deploy or Bind contract: %v", err)
	}
	err = periodicOracleUpdateHelper(numCoins, sleepSeconds, auth, contract, conn)
	if err != nil {
		log.Fatalf("failed periodic update: %v", err)
	}
	/*
	 * Update Oracle periodically with top coins
	 */
	ticker := time.NewTicker(time.Duration(frequencySeconds) * time.Second)
	go func() {
		for range ticker.C {
			err = periodicOracleUpdateHelper(numCoins, sleepSeconds, auth, contract, conn)
			if err != nil {
				log.Fatalf("failed periodic update: %v", err)
			}
		}
	}()
	select {}
}

func periodicOracleUpdateHelper(numCoins int, sleepSeconds int, auth *bind.TransactOpts, contract *diaCoingeckoOracleService.DIACoingeckoOracle, conn *ethclient.Client) error {

	topCoins, err := getTopCoinsFromCoingecko(numCoins)
	if err != nil {
		log.Fatalf("Failed to get top %d coins from Coingecko: %v", numCoins, err)
	}

	// Get quotation for topCoins and update Oracle
	for _, symbol := range topCoins {
		rawQuot, err := getForeignQuotationFromDia("Coingecko", symbol)
		if err != nil {
			log.Fatalf("Failed to retrieve Coingecko data from DIA: %v", err)
			return err
		}
		err = updateForeignQuotation(rawQuot, auth, contract, conn)
		if err != nil {
			log.Fatalf("Failed to update Coingecko Oracle: %v", err)
			return err
		}
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}

	return nil
}

func updateForeignQuotation(foreignQuotation *models.ForeignQuotation, auth *bind.TransactOpts, contract *diaCoingeckoOracleService.DIACoingeckoOracle, conn *ethclient.Client) error {
	symbol := foreignQuotation.Symbol
	price := foreignQuotation.Price
	timestamp := foreignQuotation.Time.Unix()
	err := updateOracle(conn, contract, auth, symbol, int64(price*100000), timestamp)
	if err != nil {
		log.Fatalf("Failed to update Oracle: %v", err)
		return err
	}

	return nil
}

// getTopCoinsFromCoingecko returns the symbols of the top @numCoins assets from coingecko by market cap
func getTopCoinsFromCoingecko(numCoins int) ([]string, error) {
	contents, statusCode, err := utils.GetRequest("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=" + strconv.Itoa(numCoins) + "&page=1&sparkline=false")
	if err != nil {
		return []string{}, fmt.Errorf("error on dia api with return code %d", statusCode)
	}
	type aux struct {
		Symbol string `json:"symbol"`
	}
	var quotations []aux
	err = json.Unmarshal(contents, &quotations)
	if err != nil {
		return []string{}, fmt.Errorf("error on dia api with return code %d", statusCode)
	}
	var symbols []string
	for _, quotation := range quotations {
		symbols = append(symbols, strings.ToUpper(quotation.Symbol))
	}
	return symbols, nil
}

func getForeignQuotationFromDia(source, symbol string) (*models.ForeignQuotation, error) {
	contents, _, err := utils.GetRequest(dia.BaseUrl + "v1/foreignQuotation/" + strings.Title(strings.ToLower(source)) + "/" + strings.ToUpper(symbol))
	if err != nil {
		return nil, err
	}

	var quotation models.ForeignQuotation
	err = quotation.UnmarshalBinary(contents)
	if err != nil {
		return nil, err
	}
	return &quotation, nil
}

func deployOrBindContract(deployedContract string, conn *ethclient.Client, auth *bind.TransactOpts, contract **diaCoingeckoOracleService.DIACoingeckoOracle) error {
	var err error
	if deployedContract != "" {
		*contract, err = diaCoingeckoOracleService.NewDIACoingeckoOracle(common.HexToAddress(deployedContract), conn)
		if err != nil {
			return err
		}
	} else {
		// deploy contract
		var addr common.Address
		var tx *types.Transaction
		addr, tx, *contract, err = diaCoingeckoOracleService.DeployDIACoingeckoOracle(auth, conn)
		if err != nil {
			log.Fatalf("could not deploy contract: %v", err)
			return err
		}
		log.Printf("Contract pending deploy: 0x%x\n", addr)
		log.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
		time.Sleep(180000 * time.Millisecond)
	}
	return nil
}

func updateOracle(
	client *ethclient.Client,
	contract *diaCoingeckoOracleService.DIACoingeckoOracle,
	auth *bind.TransactOpts,
	key string,
	value int64,
	timestamp int64) error {

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get 110% of the gas price
	fmt.Println(gasPrice)
	fGas := new(big.Float).SetInt(gasPrice)
	fGas.Mul(fGas, big.NewFloat(1.1))
	gasPrice, _ = fGas.Int(nil)
	fmt.Println(gasPrice)
	// Write values to smart contract
	tx, err := contract.SetValue(&bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasLimit: 800725,
		GasPrice: gasPrice,
	}, key, big.NewInt(value), big.NewInt(timestamp))
	if err != nil {
		return err
	}
	fmt.Println(tx.GasPrice())
	log.Printf("key: %s\n", key)
	log.Printf("Tx To: %s\n", tx.To().String())
	log.Printf("Tx Hash: 0x%x\n", tx.Hash())
	return nil
}
