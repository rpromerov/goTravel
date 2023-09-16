package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"

	"model"
)

func main() {
	var input string
	fmt.Println("Bienvenido a goTravel!")
	for {
		fmt.Println("1. Realizar búsqueda.")
		fmt.Println("2. Obtener reserva.")
		fmt.Println("3. Salir.")

		print("Ingrese una opción: ")
		fmt.Scanln(&input)
		switch input {
		case "1":
			{
				// Ingresar variables
				var originLocationCode, destinationLocationCode, departureDate, adults string

				// Usuario ingresa datos
				fmt.Print("Aeropuerto de origen: ")
				fmt.Scan(&originLocationCode)

				fmt.Print("Aeropuerto de destino: ")
				fmt.Scan(&destinationLocationCode)

				fmt.Print("Fecha de salida (YYYY-MM-DD): ")
				fmt.Scan(&departureDate)

				fmt.Print("Cantidad de adultos: ")
				fmt.Scan(&adults)

				// Buscar vuelos
				url := "http://localhost:5000/api/search?" + "originLocationCode=" + originLocationCode + "&destinationLocationCode=" + destinationLocationCode + "&departureDate=" + departureDate + "&adults=" + adults + "&nonStop=true&currencyCode=CLP&travelClass=ECONOMY"

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					panic(err)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				// Leer y mostrar la respuesta
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				print("Se obtuvieron los siguientes resultados:\n")

				// Crear una nueva tabla
				table := tablewriter.NewWriter(os.Stdout)
				headers := []string{"VUELO", "NÚMERO", "HORA DE SALIDA", "HORA DE LLEGADA", "AVIÓN", "PRECIO TOTAL"}
				table.SetHeader(headers)

				var flightSearchResponse model.FlightSearchResponse
				err_ := json.Unmarshal(body, &flightSearchResponse)
				if err_ != nil {
					panic(err_)
				}

				// Rellenar datos
				for _, dato := range flightSearchResponse.Data {
					for _, itinerary := range dato.Itineraries {
						for _, segment := range itinerary.Segments {

							departure := segment.Departure
							parsedTime, err := time.Parse("2006-01-02T15:04:05", departure.At)
							if err != nil {
								fmt.Println("Error al analizar la hora:", err)
								return
							}
							formattedTime := parsedTime.Format("15:04")

							arrival := segment.Arrival
							parsedArrival, err := time.Parse("2006-01-02T15:04:05", arrival.At)
							if err_ != nil {
								fmt.Println("Error al analizar la hora:", err_)
								return
							}
							formattedTimeA := parsedArrival.Format("15:04")

							row := []string{dato.ID, segment.CarrierCode + " " + segment.Number, formattedTime, formattedTimeA, "A" + segment.Aircraft.Code, dato.Price.Total}
							table.Append(row)
						}
					}
				}

				table.Render()

				// Seleccionar opción de vuelo u Otra búsqueda
				var index string
				print("Seleccione un vuelo (Ingrese 0 para realizar nueva búsqueda): ")
				fmt.Scanln(&index)
				if index != "0" {

					
				}
			}

		case "2":
			{
				var id string
				println("Ingrese el ID de la reserva: ")
				fmt.Scanln(&id)

				url := "localhost:5000/api/booking/:id" + id

				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					panic(err)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				// Leer y mostrar la respuesta
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(body))

				println("Resultado:")

				// Crear una nueva tabla
				table := tablewriter.NewWriter(os.Stdout)
				// Encabezados de la tabla
				headers := []string{"NÚMERO", "HORA DE SALIDA", "HORA DE LLEGADA", "AVIÓN", "PRECIO TOTAL"}
				// Setear los encabezados de la tabla
				table.SetHeader(headers)

				// Rellenar datos

				// Imprimir tabla
				table.Render()

				println("Pasajeros:")

				// Crear una nueva tabla
				passenger := tablewriter.NewWriter(os.Stdout)
				// Encabezados de la tabla
				headers_p := []string{"NOMBRE", "APELLIDO"}
				// Setear los encabezados de la tabla
				passenger.SetHeader(headers_p)

				// Rellenar datos

				// Imprimir tabla
				passenger.Render()
			}
		case "3":
			{
				println("Gracias por usar goTravel!")
				return
			}
		}
	}
}
