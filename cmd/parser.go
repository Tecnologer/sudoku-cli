package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var splitReg *regexp.Regexp

func init() {
	splitReg, _ = regexp.Compile(`\s`)
}
func CallCmd(cmd string, args ...string) {
	if cmd == "" {
		cmd, args = readInput("Type your next action (help for help): ")
	}

	if action, ok := commandsMap[cmd]; ok {
		action(args...)
		return
	}

	fmt.Println("Invalid option. Please, try again.")
}

func readInput(msg string) (string, []string) {
	for {
		if len(msg) > 0 {
			fmt.Println(msg)
		}
		fmt.Print(">> ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err == nil {
			input := strings.ReplaceAll(strings.Trim(input, " "), "\n", "")
			cmd := strings.SplitN(input, " ", 2)[0]
			input = strings.Trim(input[len(cmd):], " ")

			args := strings.Split(input, ",")
			return cmd, args
		}

		fmt.Println("incorrect input. try again.")
	}
}
