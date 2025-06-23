package main

import (
	"fmt"

	"github.com/kehuka11/go-yml-load-sample/config"
)

func main() {
	fmt.Printf("%+v\n", *config.Get())
}
