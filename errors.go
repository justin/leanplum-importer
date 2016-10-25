package main

import (
	"fmt"
	"os"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
		log.Error(err)
		os.Exit(1)
	}
}

func logOnError(err error) {
	if err != nil {
		log.Error(err)
	}
}
