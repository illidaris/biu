package main

import (
	"biu/pkg/property"
)

func main() {
	err := property.Input()
	if err != nil {
		println(err)
	}

}
