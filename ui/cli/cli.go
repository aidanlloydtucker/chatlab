package cli

import "fmt"

func printAll(stringChanChan <-chan chan string) {
	for {
		strChan := <-stringChanChan
		for {
			str, ok := <-strChan
			if ok {
				fmt.Printf(str)
			} else {
				break
			}
		}
		fmt.Println()
	}
}
