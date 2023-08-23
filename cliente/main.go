package main

import "fmt"

// loop waiting for user input
func main() {
	fmt.Println("Bienvenido a goTravel!")
	fmt.Println("1. Realizar búsqueda.")
	fmt.Println("2. Obtener reserva.")
	fmt.Println("3. Salir.")
	var input string
	for {
		fmt.Scanln(&input)
		switch input {
		case "1":
			println("busqueda")
			break
		case "2":
			println("obtener reserva")
			break
		case "3":
			break

		}
	}
}
