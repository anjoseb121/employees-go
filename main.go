package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrorToken  = errors.New("Error, build auth token")
	ErrorDB     = errors.New("Error, connect to DB")
	ErrorPingDB = errors.New("Error, ping to DB")
	dbEndpoint  = os.Getenv("DB_ENDPOINT")
	awsRegion   = os.Getenv("DB_REGION")
	dbUser      = os.Getenv("DB_USER")
	dbName      = os.Getenv("DB_NAME")
	// dbPassword := os.Getenv("DB_USER_PASSWORD")
)

func responseError(msg error, err error) (events.APIGatewayProxyResponse, error) {
	log.Println(err.Error())
	log.Println(msg.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       msg.Error(),
	}, msg
}

// Handler response to API Gateway requests
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	awsCreds := credentials.NewEnvCredentials()
	authToken, err := rdsutils.BuildAuthToken(dbEndpoint, awsRegion, dbUser, awsCreds)
	if err != nil {
		return responseError(ErrorToken, err)
	}
	// Create the MySQL DNS string for the DB connection
	// user:password@protocol(endpoint)/dbname?<params>
	dnsStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true",
		dbUser, authToken, dbEndpoint, dbName,
	)

	// Use db to perform SQL operations on database
	db, err := sql.Open("mysql", dnsStr)
	if err != nil {
		return responseError(ErrorDB, err)
	}

	// Test db
	err = db.Ping()
	if err != nil {
		return responseError(ErrorPingDB, err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
