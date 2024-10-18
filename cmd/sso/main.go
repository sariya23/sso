package main

import (
	"fmt"
	"sso/interanal/config"
)

func main() {
	config := config.MustLoad()
	fmt.Println(config)
}
