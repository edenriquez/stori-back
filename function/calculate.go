package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TxnBalance struct {
	Month string `json:"month"`
	Balance float64 `json:"balance"`
	TxnCounter int `json:"txn_counter"`
}

func roundToTwoDecimals(amount float64) float64 {
	rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)
	return rounded
}

func stringEquals(s1 []string, s2 string) bool {
	for _, s := range s1 {
		if s == s2 {
			return true
		}
	}
	return false
}

func groupMonths(records [][]string) []string{
	months := make([]string, 0)
	for _, record := range records {
		currentMonth := strings.Split(record[1], "/")[0]
		if !stringEquals(months, currentMonth) {
			months = append(months, currentMonth)
		}
	}
	return months
}

func sumBalances(records [][]string, month string) TxnBalance {
	var sumBalance float64
	var txnCounter int
	for _, record := range records {
		currentMonth := strings.Split(record[1], "/")[0]
		if currentMonth == month {
			balance, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				fmt.Println("Transaction balance value is incorrect:", err)
				return TxnBalance{}
			}
			fmt.Println("adding balance", balance, "for month", month)
			sumBalance += balance
			txnCounter += 1
		}
	}

	return TxnBalance{month, roundToTwoDecimals(sumBalance), txnCounter}
}


func calculateTransactionBalances(records [][]string) []TxnBalance{
	monthlyBalance := make([]TxnBalance, 0)
	months := groupMonths(records)
	for _, month := range months {
		monthlyBalance = append(
			monthlyBalance, 
			sumBalances(records, month),
		)
	}
	return monthlyBalance
}

func calculateAvgDebitAndCredit(records [][]string) (float64, float64) {
	var sumDebit float64
	var sumCredit float64
	var debitCounter int
	var creditCounter int
	for _, record := range records {
		balance, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Println("Transaction balance value is incorrect:", err)
			return 0, 0
		}
		if balance < 0 {
			sumDebit += float64(balance)
			debitCounter+=1
		} else {
			sumCredit += float64(balance)
			creditCounter+=1
		}
	}
	return sumDebit/float64(debitCounter), sumCredit/float64(creditCounter)
}
