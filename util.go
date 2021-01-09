package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/minio/sha256-simd"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getInput(query string, allowBlank bool) (input string) {
	scanner := bufio.NewScanner(os.Stdin)

	for input == "" {
		fmt.Print(query + ": ")
		scanner.Scan()
		input = strings.TrimSpace(scanner.Text())

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
