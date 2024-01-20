package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"function/db"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var S3_REGION string
var S3_BUCKET string

func init() {
	S3_REGION = os.Getenv("DEFAULT_AWS_REGION")
	S3_BUCKET = os.Getenv("S3_BUCKET")
}

const (
	customerIdParameter    string = "customerId"
	transactionIdParameter string = "transactionId"
)

func ProcessTxnFile(c *gin.Context) {
	transactionId := uuid.New().String()
	fileHeader, err := c.FormFile("file")
	customerId := c.Param(customerIdParameter)
	email := c.PostForm("email")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sourceFile io.ReadCloser
	sourceFile, err = fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer sourceFile.Close()

	sess := session.Must(session.NewSession())
	sess.Config.WithRegion(S3_REGION)

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})

	now := time.Now()

	filePath := fmt.Sprintf("%s/transactions/%s/%s.csv", customerId,  transactionId, now.Format("2006-01-02"))

	upResponse, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(S3_BUCKET),
		Key:    aws.String(filePath),
		Body:   sourceFile,

	})
	if err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("file uploaded succesfuly...")
	fmt.Println(upResponse.Location)
	

	fileLocation := upResponse.Location

	reader := csv.NewReader(sourceFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	balances := calculateTransactionBalances(records[1:])
	avgDebit, avgCredit := calculateAvgDebitAndCredit(records[1:])
	totalBalance := float64(0)

	for _, balance := range balances {
		totalBalance += float64(balance.Balance)
	}

	sent, err := send_email(
		email,
		totalBalance,
		balances,
		avgDebit,
		avgCredit,
	)

	fmt.Println("email was sent to", email)
	fmt.Println(sent)
	if err!= nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balancesMap := make([]map[string]interface{}, len(balances))

	for index, b := range balances {
		balancesMap[index] = map[string]interface{}{
			"month": b.Month,
			"balance": b.Balance,
			"txnCounter": b.TxnCounter,
		}
	}

	db.SaveTransaction(
		roundToTwoDecimals(totalBalance),
		roundToTwoDecimals(avgDebit),
		roundToTwoDecimals(avgCredit),
		balancesMap,
		transactionId,
		customerId,
		fileLocation,
	)

	c.JSON(http.StatusAccepted, gin.H{
		"transaction_id":        transactionId,
		"balances":              balances,
		"total_balance":         roundToTwoDecimals(totalBalance),
		"average_debit_amount":  roundToTwoDecimals(avgDebit),
		"average_credit_amount": roundToTwoDecimals(avgCredit),
		"file":                  fileLocation,
		"success":               true,
	})
}

func GetTransaction(c *gin.Context){
	fmt.Println("Getting transaction details...")
	
	transactionId := c.Param(transactionIdParameter)
	summary, err := db.GetTransactionDetails(transactionId)
	
	if err!= nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	
	var balances []map[string]interface{}
	json.Unmarshal([]byte(summary["balances"].(string)), &balances)

	c.JSON(http.StatusAccepted, gin.H{
		"transaction_id":        summary["transactionId"],
		"total_balance":         summary["totalBalance"],
		"balances":              balances,
		"average_debit_amount":  summary["avgDebit"],
		"average_credit_amount": summary["avgCredit"],
		"file":                  summary["file"],
		"success":               true,
	})

}

func GetCustomerTransactions(c *gin.Context){
	fmt.Println("Getting customer transactions...")
	
	customerId := c.Param(customerIdParameter)
	transactions, err := db.GetTransactions(customerId)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var arrSummaries []map[string]interface{}
	
	for _, transaction := range transactions {
		var balances []map[string]interface{}
		json.Unmarshal([]byte(transaction["balances"].(string)), &balances)
		arrSummaries = append(arrSummaries, map[string]interface{}{
			"transaction_id":        transaction["transactionId"],
			"total_balance":         transaction["totalBalance"],
			"balances":              balances,
			"average_debit_amount":  transaction["avgDebit"],
			"average_credit_amount": transaction["avgCredit"],
			"file":                  transaction["file"],
		})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"transactions": arrSummaries,
		"success": true, 
	})

}


