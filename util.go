package main

import (
	"fmt"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getInput(query string, allowBlank bool) (input string) {
	for input == "" {
		fmt.Print(query + ": ")
		fmt.Scanf("%s\n", &input)

		if allowBlank {
			break
		}
	}

	return
}
