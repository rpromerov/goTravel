package cliente

import "fmt"

// loop waiting for user input
func main() {
	var input string
	for {
		fmt.Scanln(&input)
		if input == "exit" {
			break
		}
	}
}
