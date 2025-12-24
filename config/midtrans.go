package config

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewMidtransClient() coreapi.Client {
	var client coreapi.Client

	env := midtrans.Sandbox
	if os.Getenv("MIDTRANS_ENV") == "production" {
		env = midtrans.Production
	}

	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), env)
	return client
}
