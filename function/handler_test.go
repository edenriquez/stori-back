package main

import (
	"encoding/csv"
	"io"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)
var records [][]string

func init() {
	var sourceFile io.ReadCloser
	file, _ := os.Open("./examples/transactions_min.csv")
	sourceFile = file

	defer sourceFile.Close()

	reader := csv.NewReader(sourceFile)
	records, _ = reader.ReadAll()
}

func TestCalculateTransactionBalances(t *testing.T) {
	balances := calculateTransactionBalances(records[1:])
	assert.Equal(t, TxnBalance{"7", 50.2, 2}, balances[0])
}
