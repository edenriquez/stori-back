package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var client *dynamodb.Client


const customerIdDynamoDBAttributeName = "customerId"
const transactionIdDynamoDBAttributeName = "transactionId"
const avgDebitDynamoDBAttributeName = "avgDebit"
const avgCreditDynamoDBAttributeName = "avgCredit"
const totalBalanceDynamoDBAttributeName = "totalBalance"
const fileLocationDynamoDBAttributeName = "file"
const balancesDynamoDBAttributeName = "balances"

var ErrTxnNotFound = errors.New("txn not found")

var table string
var region string

func init() {
	table = os.Getenv("TABLE_NAME")
	region = os.Getenv("DEFAULT_AWS_REGION")
	fmt.Println("initializing ddb client for table", table)

	cfg, _ := config.LoadDefaultConfig(context.Background())
	cfg.Region = region
	client = dynamodb.NewFromConfig(cfg)

	fmt.Println("db client initialized...")
}

func SaveTransaction(
		totalBalance, 
		avgDebit, 
		avgCredit float64, 
		balances []map[string]interface{}, 
		transactionId, 
		customerId string,
		fileLocation string,
	) (string, error) {
	item := make(map[string]types.AttributeValue)

	item[transactionIdDynamoDBAttributeName] = &types.AttributeValueMemberS{Value: transactionId}
	item[customerIdDynamoDBAttributeName] = &types.AttributeValueMemberS{Value: customerId}
	item[avgDebitDynamoDBAttributeName] = &types.AttributeValueMemberN{Value: strconv.FormatFloat(avgDebit, 'f', -1, 64)}
	item[avgCreditDynamoDBAttributeName] = &types.AttributeValueMemberN{Value: strconv.FormatFloat(avgCredit, 'f', -1, 64)}
	item[totalBalanceDynamoDBAttributeName] = &types.AttributeValueMemberN{Value: strconv.FormatFloat(totalBalance, 'f', -1, 64)}
	item[fileLocationDynamoDBAttributeName] = &types.AttributeValueMemberS{Value: fileLocation}

	balancesJSON, err := json.Marshal(balances)
	if err != nil {
		return "", err
	}
	item[balancesDynamoDBAttributeName] = &types.AttributeValueMemberS{Value: string(balancesJSON)}
	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      item})

	if err != nil {
		log.Println("dynamodb put item failed", err)
		return "", err
	}

	log.Printf("inserting transaction: %s for customer: %s\n", transactionId, customerId)
	return transactionId, nil
}



func GetTransactions(customerId string) ([]map[string]interface{}, error) {
	fmt.Println(customerId)
	
	txns, err := client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(table),
		IndexName:              aws.String("customerId"),
		KeyConditionExpression: aws.String("customerId = :customerId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":customerId": &types.AttributeValueMemberS{Value: customerId},
		},
	})

	if err != nil {
		log.Println("failed to get transaction details", err)
		return nil, err
	}

	if txns.Items == nil {
		return nil, ErrTxnNotFound
	}

	var transactions []map[string]interface{}

	for _, item := range txns.Items {
		var itemMap map[string]interface{}
		attributevalue.UnmarshalMap(item, &itemMap)
		transactions = append(transactions, itemMap)
	}

	return transactions, nil
}

func GetTransactionDetails(transactionId string) (map[string]interface{}, error) {
	fmt.Println(transactionId)
	txn, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			transactionIdDynamoDBAttributeName: &types.AttributeValueMemberS{Value: transactionId},
		},
	})

	if err != nil {
		log.Println("failed to get transactions", err)
		return nil, err
	}

	if txn.Item == nil {
		return nil, ErrTxnNotFound
	}

	var summary map[string]interface{}

	attributevalue.UnmarshalMap(txn.Item, &summary)

	return summary, nil
}
