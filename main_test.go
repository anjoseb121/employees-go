package main

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	os.Setenv("DB_ENDPOINT", "mysqllambda.c2vw70w2ufaf.us-east-2.rds.amazonaws.com")
	os.Setenv("DB_REGION", "us-east-2a")
	os.Setenv("DB_USER", "username")
	os.Setenv("DB_NAME", "ExampleDB")

	tests := []struct {
		request events.APIGatewayProxyRequest
		err     error
	}{
		{
			// Test the ping to db
			request: events.APIGatewayProxyRequest{},
			err:     nil,
		},
	}

	for _, test := range tests {
		_, err := Handler(test.request)
		assert.IsType(t, test.err, err)
	}
}
