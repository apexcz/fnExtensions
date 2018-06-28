package main

import(
	"context"
	"fmt"
	"github.com/fnproject/fn/api/server"
	_ "github.com/apexcz/fnExtensions/smartlog"
)

func main() {
	fmt.Println("Omoh")
	ctx := context.Background()
	funcServer := server.NewFromEnv(ctx)
	funcServer.AddExtensionByName("github.com/apexcz/fnExtensions/smartlog")
	funcServer.Start(ctx)
}