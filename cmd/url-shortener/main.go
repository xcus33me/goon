package main

import (
	"fmt"
	"glearn/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
