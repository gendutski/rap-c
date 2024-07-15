package config

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// generate .env file from struct envconfig tags
func GenerateDotEnv() {
	cfg := config{}
	rows := readStruct(reflect.TypeOf(cfg))
	text := strings.Join(rows, "\n")
	err := os.WriteFile(".env", []byte(text), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readStruct(elm reflect.Type) []string {
	var result []string
	reader := bufio.NewReader(os.Stdin)
	numFields := elm.NumField()
	for i := 0; i < numFields; i++ {
		// set field, tag, kind
		field := elm.Field(i)
		tag := field.Tag
		kind := field.Type.Kind()

		// env variable
		envconfig := tag.Get("envconfig")
		if envconfig == "" || envconfig == "-" {
			continue
		}

		// init prompt
		prompt := tag.Get("prompt")
		_default := tag.Get("default")
		if _default != "" {
			prompt += fmt.Sprintf(" (default:%s)", _default)
		}
		if kind.String() == "bool" {
			prompt += " (true or false)"
		}
		prompt += ": "

		// is secret prompt
		secret, _ := strconv.ParseBool(tag.Get("secret"))

		// scan
		var scan string
		if secret {
			scan = promptPassword(prompt)
			fmt.Println()
		} else {
			scan = promptString(prompt, reader)
		}
		if scan == "" && _default != "" {
			scan = _default
		}
		result = append(result, fmt.Sprintf("%s=\"%s\"", envconfig, strings.ReplaceAll(scan, `"`, "")))
	}
	return result
}

func promptString(prompt string, reader *bufio.Reader) string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	fmt.Print(prompt)
	result, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error getString:", err.Error())
		os.Exit(1)
	}
	result = strings.TrimSpace(result)
	return result
}

func promptPassword(prompt string) string {
	fmt.Print(prompt)
	result, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Println("error getString:", err.Error())
		os.Exit(1)
	}
	return strings.TrimSpace(string(result))
}
