package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const functionDir = "../function"

type LambdaStoriBackendProps struct {
	awscdk.StackProps
}

func InfraestructureCreation(scope constructs.Construct, id string, props *LambdaStoriBackendProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Lambda infrastructure creation

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("gin-go-lambda-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime: awslambda.Runtime_GO_1_X(),
			Entry:   jsii.String(functionDir),
		})

	// Api gateway infrastructure creation

	api := awsapigateway.NewLambdaRestApi(stack, jsii.String("lambda-rest-api"), &awsapigateway.LambdaRestApiProps{
		Handler: function,
	})

	app := api.Root().AddResource(jsii.String("api"), nil)

	// Api gateway endpoints declaration

	customer := app.AddResource(jsii.String("{customerId}"), nil)
	transactions := customer.AddResource(jsii.String("transactions"), nil)
	TxnDetails := transactions.AddResource(jsii.String("{transactionId}"), nil)

	
	
	
	transactions.AddMethod(jsii.String("POST"), nil, nil) // GET /api/{customerId}/transactions
	transactions.AddMethod(jsii.String("GET"), nil, nil)  // GET /api/{customerId}/transactions
	TxnDetails.AddMethod(jsii.String("GET"), nil, nil)   // GET /api/{customerId}/transactions/{transactionId}

	awscdk.NewCfnOutput(stack, jsii.String("api-gateway-endpoint"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("API-Gateway-Endpoint"),
			Value:      api.Url()})

	// S3 bucket creation

	s3Bucket := awss3.NewBucket(stack, jsii.String("StoriStorageTransactionFiles"), &awss3.BucketProps{
		BucketName: jsii.String("stori-storage-transaction-file"),
	})

	// Granting S3 permissions to Lambda function
	
	s3Bucket.GrantReadWrite(function, nil)
	

	// create dynamo postgres database
	table := awsdynamodb.NewTable(stack, jsii.String("StoriTansactionsInformation"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("transactionId"),
				Type: awsdynamodb.AttributeType_STRING,
			},
		})
	secondaryIndex := awsdynamodb.GlobalSecondaryIndexProps{
			IndexName: jsii.String("customerId"),
			ProjectionType: awsdynamodb.ProjectionType_ALL,
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("customerId"),
				Type: awsdynamodb.AttributeType_STRING,
			},
		}
	table.AddGlobalSecondaryIndex(&secondaryIndex)
	table.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)
	table.GrantReadWriteData(function)

	awscdk.NewCfnOutput(stack, jsii.String("stori-transaction-table"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("stori-transaction-table"),
			Value:      table.TableName()})

	return stack
}



func main() {
	app := awscdk.NewApp(nil)

	InfraestructureCreation(app, "LambdaStoriBackend", &LambdaStoriBackendProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
