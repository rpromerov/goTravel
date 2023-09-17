package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"

	"model"
)

func main() {
	fmt.Println("Bienvenido a goTravel!")
	for {
		var input string
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
					fmt.Println("Hubo un error al realizar la búsqueda”", err)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Hubo un error al realizar la búsqueda”", err)
				}
				defer resp.Body.Close()

				// Leer y mostrar la respuesta
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Hubo un error al realizar la búsqueda”", err)
				}

				print("Se obtuvieron los siguientes resultados:\n")

				// Crear una nueva tabla
				table := tablewriter.NewWriter(os.Stdout)
				headers := []string{"VUELO", "NÚMERO", "HORA DE SALIDA", "HORA DE LLEGADA", "AVIÓN", "PRECIO TOTAL"}
				table.SetHeader(headers)

				var flightSearchResponse model.FlightSearchResponse
				err_ := json.Unmarshal(body, &flightSearchResponse)
				if err_ != nil {
					fmt.Println("Hubo un error al realizar la búsqueda”", err)
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

							row := []string{dato.ID, segment.CarrierCode + segment.Number, formattedTime, formattedTimeA, "A" + segment.Aircraft.Code, dato.Price.Total}
							table.Append(row)
						}
					}
				}

				table.Render()

				// Seleccionar opción de vuelo u Otra búsqueda
				var index string
				fmt.Print("Seleccione un vuelo (Ingrese 0 para realizar nueva búsqueda): ")
				fmt.Scan(&index)

				var jsonData model.FlightData

				if index != "0" {
					for _, dato := range flightSearchResponse.Data {
						if dato.ID == index {
							dataInstance := model.Data{
								Type:         "flight-offers-pricing",
								FlightOffers: []model.FlightOffer{dato},
							}
							jsonData = model.FlightData{
								Data: dataInstance,
							}
							break

						}
					}

				}

				jsonDataBytes, err := json.Marshal(jsonData)
				if err != nil {
					panic(err)
				}

				URL := "http://localhost:5000/api/pricing"

				jsonDataReader := bytes.NewReader(jsonDataBytes)

				// Crear una nueva solicitud POST con jsonDataReader como cuerpo de la solicitud
				sol, err := http.NewRequest("POST", URL, jsonDataReader)
				if err != nil {
					panic(err)
				}
				sol.Header.Set("Content-Type", "application/json")

				// Realizar la solicitud HTTP
				cliente := &http.Client{}
				ans, err := cliente.Do(sol)
				if err != nil {
					panic(err)
				}
				defer ans.Body.Close()

				responseBody, err := ioutil.ReadAll(ans.Body)
				if err != nil {
					panic(err)
				}

				var response model.FlightOfferResponse
				err2 := json.Unmarshal([]byte(responseBody), &response)
				if err2 != nil {
					panic(err)
				}

				Total := response.Data.FlightOffers[0].Price.Total

				fmt.Println("El precio total final es de: ", Total)

				// Solicita la cantidad de pasajeros
				adultsAsInt, err := strconv.Atoi(adults)
				if err != nil {
					fmt.Println("Error al convertir el string a número:", err)
					return
				}

				for i := 1; i <= adultsAsInt; i++ {
					fmt.Println("Pasajero", i)

					var born, name, lastname, gender, mail, phone string

					fmt.Print("Ingrese fecha de nacimiento: ")
					_, _ = fmt.Scan(&born)

					fmt.Print("Ingrese nombre: ")
					_, _ = fmt.Scan(&name)

					fmt.Print("Ingrese apellido: ")
					_, _ = fmt.Scan(&lastname)

					fmt.Print("Ingrese sexo (MALE o FEMALE): ")
					_, _ = fmt.Scan(&gender)

					fmt.Print("Ingrese correo: ")
					_, _ = fmt.Scan(&mail)

					fmt.Print("Ingrese teléfono: ")
					_, _ = fmt.Scan(&phone)

					// Insertar la reserva en MongoDB (tendrás que definir la estructura `Contact` y `TravelerInfo`)
					contact := model.Contact{
						EmailAddress: mail,
						Phones: []model.Phone{
							{Number: phone},
						},
					}

					traveler := model.TravelerInfo{
						ID:          uuid.New().String(),
						DateOfBirth: born,
						Name: model.Name{
							FirstName: name,
							LastName:  lastname,
						},
						Gender:  gender,
						Contact: contact,
					}

					// Generar un ID único
					id := uuid.New().String()

					// Crear una nueva instancia de FlightOrderData
					flightOrder := model.FlightOrderData{
						Type:            "flight-order",
						ID:              id,
						QueuingOfficeId: "Amadeus123",
						Travelers:       []model.TravelerInfo{traveler},
					}

					// Convierte los datos en JSON
					requestBody, err := json.Marshal(flightOrder)
					if err != nil {
						fmt.Println("Error al convertir a JSON:", err)
						return
					}

					// URL de la API donde deseas hacer la solicitud POST
					apiURL := "http://localhost:5000/api/booking" // Ajusta la URL a tu entorno

					// Realiza la solicitud POST
					resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
					if err != nil {
						fmt.Println("Error al hacer la solicitud POST:", err)
						return
					}
					defer resp.Body.Close()

					// Lee la respuesta del servidor si es necesario
					responseBody, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("Error al leer la respuesta del servidor:", err)
						return
					}

					fmt.Println("Respuesta del servidor:", string(responseBody))

					fmt.Println("Reserva insertada exitosamente: " + id)
				}
			}

		case "2":
			{
				var id string
				fmt.Print("Ingrese el ID de la reserva: ")
				fmt.Scan(&id)

				url := fmt.Sprintf("http://localhost:5000/api/booking/%s", id)

				// Realiza la solicitud GET al endpoint
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("Hubo un error al recuperar la reserva:", err)
					return
				}
				defer resp.Body.Close()

				// Verifica si la respuesta es exitosa (código de estado 200)
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("Hubo un error al recuperar la reserva: %d\n", resp.StatusCode)
					return
				}

				// Lee y decodifica la respuesta JSON
				var response model.FlightOrderResponse
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					fmt.Println("Hubo un error al recuperar la reserva", err)
					return
				}
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
