package main

import (
	"fmt"

	"github.com/fusor/cpma/cmd"
	"github.com/fusor/cpma/env"
)

func main() {
	cmd.Execute()
	config := env.New()
	config.Fetch()
	fmt.Println(config)
}
