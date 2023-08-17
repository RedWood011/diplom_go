package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"RedWood011/client/adapter/secret"
	"RedWood011/client/adapter/user"
	"RedWood011/client/cli"
	secretservice "RedWood011/client/service/secret"
	userservice "RedWood011/client/service/user"
)

var (
	//nolint: gochecknoglobals //Flag initialization
	BuildTime string
	//nolint: gochecknoglobals //Flag initialization
	AppVersion string
)

const defaultAddress = "localhost:5050"

func main() {
	ctx := context.Background()
	address := flag.String("a", defaultAddress, "address of gGRPC server")
	if *address == defaultAddress {
		envAddress := os.Getenv("ADDRESS")
		if envAddress != "" {
			address = &envAddress
		}
	}

	fmt.Println("InitApp")
	fmt.Printf("App version: %v, Date compile: %v\n", AppVersion, BuildTime)
	userAdapter := user.NewUserAdapter(*address)
	userService := userservice.NewUserService(userAdapter)
	secretAdapter := secret.NewSecretAdapter(*address)
	secretService := secretservice.NewSecretService(secretAdapter)

	cli.RunCli(ctx, userService, secretService)
}
