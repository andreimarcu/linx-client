package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
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

func validateCommand(cmdAndArgs string) error {
	split, err := shellquote.Split(cmdAndArgs)
	if err != nil {
		return err
	}

	if len(split) == 0 {
		return fmt.Errorf("No command supplied")
	}

	cmd := split[0]
	if _, err = exec.LookPath(cmd); err != nil {
		return fmt.Errorf("Command not found: %s", cmd)
	}

	return nil
}

func splitCmdAndArgs(cmdAndArgs string) (string, []string, error) {
	split, err := shellquote.Split(cmdAndArgs)
	return split[0], split[1:], err
}

func runCmdFirstLine(cmdAndArgs string) (string, error) {
	cmd := exec.Command(cmdAndArgs)
	var stdout io.ReadCloser
	var err error

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return "", fmt.Errorf("read from stdout: %w", err)
	}

	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("start program: %s: %w", cmdAndArgs, err)
	}

	defer stdout.Close()
	scanner := bufio.NewReader(stdout)

	line, err := scanner.ReadString('\n')

	if err == io.EOF {
		return line, nil
	}

	return line, err
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
