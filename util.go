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
	_, _, err := splitCmdAndArgs(cmdAndArgs)
	return err
}

func splitCmdAndArgs(cmdAndArgs string) (string, []string, error) {
	split, err := shellquote.Split(cmdAndArgs)

	if err != nil {
		return "", nil, err
	}

	if len(split) == 0 {
		return "", nil, fmt.Errorf("no command supplied")
	}

	cmd := split[0]
	if _, err = exec.LookPath(cmd); err != nil {
		return "", nil, fmt.Errorf("command not found: %s", cmd)
	}

	return split[0], split[1:], err
}

func runCmdFirstLine(cmdAndArgs string) (string, error) {
	path, args, err := splitCmdAndArgs(cmdAndArgs)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return "", fmt.Errorf("read from stdout: %w", err)
	}

	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("start program: %s: %w", cmdAndArgs, err)
	}

	scanner := bufio.NewReader(stdout)

	line, err := scanner.ReadString('\n')

	if err := cmd.Wait(); err != nil {
		return line, err
	}

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
