package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type config struct {
	Next     string
	Previous string
}

func main() {
	c := &config{}
	_ = GlobalCache()
	scanner := bufio.NewScanner(os.Stdin)
	//scanner.Split(bufio.ScanLines)
	for {
		fmt.Print("Pokedex > ")
		scanOk := scanner.Scan()
		if !scanOk {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			} else {
				fmt.Println("\nEOF detected. Exiting.")
			}
			break
		}
		enteredText := scanner.Text()
		if enteredText == "" {
			continue
		}
		enteredText = strings.ToLower(enteredText)
		fields := strings.Fields(enteredText)
		command := fields[0]
		var params []string
		if len(fields) > 1 {
			params = fields[1:]

		}
		cliC, ok := registry[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		} else {
			err := cliC.callback(c, params...)
			if err != nil {
				fmt.Println(err.Error())
				//os.Exit(1)
				continue
			}
		}

	}
}

func cleanInput(text string) []string {
	splitTextSlice := strings.Fields(text)
	for i, w := range splitTextSlice {
		splitTextSlice[i] = strings.ToLower(w)
	}
	return splitTextSlice
}
