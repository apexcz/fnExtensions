package main

import(
	"context"
	"fmt"
	_ "github.com/fnproject/fn/api/server/defaultexts"
	_ "github.com/fnproject/fn/api/models"

	"github.com/fnproject/fn/api/server"
	_ "github.com/apexcz/fnExtensions/smartlog"
)

func main() {
	fmt.Println("It works")
	ctx := context.Background()
	funcServer := server.NewFromEnv(ctx)
	funcServer.AddExtensionByName("github.com/apexcz/fnExtensions/smartlog")
	funcServer.Start(ctx)
}