package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrorToken  = errors.New("Error, build auth token")
	ErrorDB     = errors.New("Error, connect to DB")
	ErrorPingDB = errors.New("Error, ping to DB")
	ErrorRead   = errors.New("Error reading employees table")
)

type Employee struct {
	Id   int
	Name string
}

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
	DbEndpoint := os.Getenv("DB_ENDPOINT")
	AwsRegion := os.Getenv("DB_REGION")
	DbUser := os.Getenv("DB_USER")
	DbName := os.Getenv("DB_NAME")

	awsCreds := credentials.NewEnvCredentials()
	authToken, err := rdsutils.BuildAuthToken(DbEndpoint, AwsRegion, DbUser, awsCreds)
	if err != nil {
		return responseError(ErrorToken, err)
	}
	// Create the MySQL DNS string for the DB connection
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
		DbUser, authToken, DbEndpoint, DbName,
	)

	// Use db to perform SQL operations on database
	db, err := sql.Open("mysql", connectStr)
	defer db.Close()
	if err != nil {
		return responseError(ErrorDB, err)
	}

	// Test db
	err = db.Ping()
	if err != nil {
		return responseError(ErrorPingDB, err)
	}

	// Create table
	//db.Exec("create table Employee3 ( EmpID  int NOT NULL, Name varchar(255) NOT NULL, PRIMARY KEY (EmpID))")

	//db.Exec("insert into Employee3 (EmpID, Name) values(1, \"Joe\")")
	//db.Exec("insert into Employee3 (EmpID, Name) values(2, \"Bob\")")
	//db.Exec("insert into Employee3 (EmpID, Name) values(3, \"Mary\")")

	rows, err := db.Query("select * from Employee3")
	defer rows.Close()

	for rows.Next() {
		employee := Employee{}
		err = rows.Scan(&employee.Id, &employee.Name)
		if err != nil {
			return responseError(ErrorRead, err)
		}
		log.Printf("Employee name: %s and ID: %s", employee.Name, employee.Id)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
