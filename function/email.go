package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"function/templates"

	"github.com/resend/resend-go/v2"
)


func send_email(to string, totalBalance float64, balances []TxnBalance, avgDebit float64, avgCredit float64) (*resend.SendEmailResponse, error) {
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	txnDetailsTemplate := ""
	for _, balance := range balances {
		m, _ := strconv.Atoi(balance.Month)
		month, _ := getMonthName(m)
		txnDetailsTemplate += "<p>Number of transactions in " + month + " :" + strconv.Itoa(balance.TxnCounter) + "</p>"
	}

	template := string(templates.TEMPLATE_STRING)
	template = strings.Replace(template, "{{totalBalance}}", strconv.FormatFloat(totalBalance, 'f', -1, 64), -1)
	template = strings.Replace(template, "{{transactionDetail}}", txnDetailsTemplate, -1)
	template = strings.Replace(template, "{{avgDebit}}", strconv.FormatFloat(math.Abs(avgDebit), 'f', -1, 64), -1)
	template = strings.Replace(template, "{{avgCredit}}", strconv.FormatFloat(avgCredit, 'f', -1, 64), -1)


	params := &resend.SendEmailRequest{
	    From:    "onboarding@resend.dev",
	    To:      []string{to},
		Cc:      []string{ "itenriquez.isc@gmail.com"},
	    Subject: "Transaction Report 💸 from Stori",
	    Html:    template,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
	    fmt.Println(err.Error())
	    return nil, err
	}

	return sent, nil
}


