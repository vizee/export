package main

import (
	"os"

	"github.com/vizee/export"
	"github.com/vizee/export/linehandler"
)

var (
	testInt int = 1
)

func main() {
	export.Int("test", &testInt)

	err := linehandler.Serve(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}
