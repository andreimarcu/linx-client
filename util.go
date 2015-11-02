package main

import (
	"crypto/sha256"
	"encoding/hex"

	"fmt"
	"io"
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

func sha256sum(r io.Reader) string {
	hasher := sha256.New()

	_, err := io.Copy(hasher, r)
	if err == nil {
		return hex.EncodeToString(hasher.Sum(nil))
	} else {
		return ""
	}
}
