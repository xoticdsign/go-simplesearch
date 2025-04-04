package main

import (
	"os"

	"github.com/xoticdsign/simplesearch/internal/app"
)

/*

THIS PROJECT WAS CREATED AS AN EXAMPLE IMPLEMENTATION OF ELASTICSEARCH IN GOLANG

*/

func main() {
	var env string

	env = os.Getenv("env")
	if env == "" {
		env = "local"
	}

	a, err := app.New(env)
	if err != nil {
		panic(err)
	}
	err = a.Run()
	if err != nil {
		panic(err)
	}
}
