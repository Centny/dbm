package main

import (
	"fmt"
	"os"
)

func main() {
	// tmgo()
	if len(os.Args) < 2 {
		fmt.Println("Usage:dbmtest <tsql|tmgo>")
		return
	}
	if os.Args[1] == "tsql" {
		tsql()
	} else {
		tmgo()
	}
}
