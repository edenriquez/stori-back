package main

import (
	"fmt"
	"os"
	"time"
)

func isLocal() bool {
	return os.Getenv("ENV") == "development"
}

func getMonthName(monthNumber int) (string, error) {
	if monthNumber < 1 || monthNumber > 12 {
		return "", fmt.Errorf("invalid month number: %d", monthNumber)
	}

	month := time.Date(2022, time.Month(monthNumber), 1, 0, 0, 0, 0, time.UTC)

	return month.Format("January"), nil
}
