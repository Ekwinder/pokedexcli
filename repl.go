package main

import (
	"bufio"
	"fmt"
	"os"
)

func startRepl() {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		input.Scan()
		cmd := cleanInput(input.Text())
		if len(cmd) == 0 {
			continue
		}
		// we might want to run the command now?
		v, ok := getCommands()[cmd[0]]
		if ok {
			v.callback()
		} else {
			fmt.Println("Error: Invalid Command. Use help to see the command list.")
		}

	}
}
