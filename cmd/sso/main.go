package main

import (
	"fmt"
	"github.com/finlleyl/gRPC/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
