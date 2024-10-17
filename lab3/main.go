package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

const (
	incomeTax        = 0.2
	expensesAndTaxes = 500.0
	rentPayment      = 200.0
	transferPayment  = 150.0
	partyBaseVolume  = 200.0
	meanPrice        = 100.0
	maxDemand        = 30.0
	initAccount      = 10000.0
	commonOffer      = 50.0
	basePrice        = 35.0
)

var (
	chosenPrice         = 100.0
	shopStorage         = 30.0
	currentTime         = 0.0
	chosenCreditPayment = 0.0
	chosenCreditPeriod  = 0.0
	basicStorage        = 80.0
	account             = initAccount
)

func main() {
	for account > 0 {
		input := ""

		fmt.Printf("Account: %v\n", account)

		partyVolume := partyBaseVolume * (0.75 + rand.Float64()*(1.25-0.75))
		fmt.Printf("Party volume: %v\n", partyVolume)

		addPriceByTime := basePrice + 0.03*currentTime + basePrice*0.01*currentTime*rand.Float64()
		basicPrice := commonOffer * (0.7 + rand.Float64()*(1.3-0.7))
		unitPrice := addPriceByTime + basicPrice
		fmt.Printf("Unit price: %v\n", unitPrice)

		partyPrice := unitPrice * partyVolume
		fmt.Printf("Party price: %v\n", partyPrice)

		fmt.Print("By party?: ")
		fmt.Scanln(&input)

		if input == "yes" || input == "y" {
			if account >= partyPrice {
				account -= partyPrice
				basicStorage += partyVolume
			} else {
				fmt.Printf("Not enough money for party needed %v got %v\n", partyPrice, account)
			}
		}

		fmt.Printf("Shop storage: %v\n", shopStorage)
		fmt.Printf("Base storage: %v\n", basicStorage)

		fmt.Print("Move from base to shop?: ")
		fmt.Scanln(&input)

		if input == "yes" || input == "y" {
			fmt.Print("How many?: ")
			fmt.Scanln(&input)
			toTransfer, _ := strconv.ParseFloat(input, 64)

			if account >= transferPayment {
				toTransfer = math.Min(basicStorage, toTransfer)

				basicStorage -= toTransfer

				shopStorage += toTransfer
				if shopStorage > 100.0 {
					shopStorage = 100.0
				}

				account -= transferPayment
			} else {
				fmt.Printf("Not enough money for transfer needed %v got %v\n", transferPayment, account)
			}
		}

		demand := maxDemand * (1 - (1 / (1 + math.Exp(-0.05*(chosenPrice-meanPrice))))) * (0.7 + rand.Float64()*(1.2-0.7))
		fmt.Printf("Current demand: %v\n", demand)

		fmt.Print("Stop sells?: ")
		fmt.Scanln(&input)

		if input == "no" || input == "n" {
			fmt.Printf("Choose new price (old - %v):", chosenPrice)
			fmt.Scanln(&input)
			chosenPrice, _ = strconv.ParseFloat(input, 64)

			sells := math.Min(demand, shopStorage)
			shopStorage -= sells

			income := chosenPrice * sells

			account += income
			account -= income * incomeTax
		}

		creditValue := float64(rand.Intn(20000)) + 500
		fmt.Printf("Credit value: %v\n", creditValue)

		creditRate := rand.Float64()
		fmt.Printf("Credit rate: %v\n", creditRate)

		creditPeriod := float64(rand.Intn(10)) + 1
		fmt.Printf("Credit period: %v\n", creditPeriod)

		fmt.Print("Getting credit?: ")
		fmt.Scanln(&input)

		if input == "yes" || input == "y" {
			chosenCreditPeriod = creditPeriod
			chosenCreditPayment = creditValue * math.Pow(1+creditRate, chosenCreditPeriod) / creditPeriod
			account += creditValue
		}

		if chosenCreditPeriod > 0 {
			account -= chosenCreditPayment
			chosenCreditPeriod--
		}

		account -= math.Min(rentPayment+expensesAndTaxes, account)

		currentTime += 1
	}
}
