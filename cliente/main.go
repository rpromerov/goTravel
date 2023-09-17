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

	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"

	"model"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error al cargar las variables de entorno:", err)
		return
	}
	url := os.Getenv("SERVER")
	port := os.Getenv("PORT")
	connection := url + ":" + port
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
				url := "http://" + connection + "/api/search?" + "originLocationCode=" + originLocationCode + "&destinationLocationCode=" + destinationLocationCode + "&departureDate=" + departureDate + "&adults=" + adults + "&nonStop=true&currencyCode=CLP&travelClass=ECONOMY"

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

				} else {
					continue
				}

				jsonDataBytes, err := json.Marshal(jsonData)
				if err != nil {
					panic(err)
				}

				URL := "http://" + connection + "/api/pricing"

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
				travelers := []model.TravelerInfo{}

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
							{Number: phone,
								CountryCallingCode: "56",
								DeviceType:         "MOBILE"},
						},
					}

					traveler := model.TravelerInfo{
						ID:          response.Data.FlightOffers[0].TravelerPricings[i-1].TravelerID,
						DateOfBirth: born,
						Name: model.Name{
							FirstName: name,
							LastName:  lastname,
						},
						Gender:  gender,
						Contact: contact,
					}
					travelers = append(travelers, traveler)
				}
				var offers []model.FlightOffer

				for _, offerData := range flightSearchResponse.Data {
					// Crea una variable model.FlightOffer para cada elemento en flightSearchResponse.Data
					offer := model.FlightOffer{
						Type:                     offerData.Type,
						ID:                       offerData.ID,
						Source:                   offerData.Source,
						InstantTicketingRequired: offerData.InstantTicketingRequired,
						NonHomogeneous:           offerData.NonHomogeneous,
						// Continúa asignando los demás campos según sea necesario
					}

					// Agrega cada oferta a la lista 'offers'
					offers = append(offers, offer)
				}

				// Crear una nueva instancia de FlightOrderData
				flightOrder := model.FlightOrderData{
					Type:            "flight-order",
					QueuingOfficeId: "Amadeus123",
					FlightOffers:    []model.FlightOffer{response.Data.FlightOffers[0]},
					Travelers:       travelers,
				}
				flightOrderResponse := model.FlightOrderResponse{
					Data: flightOrder,
				}

				// Convierte los datos en JSON
				requestBody, err := json.Marshal(flightOrderResponse)
				if err != nil {
					fmt.Println("Error al convertir a JSON:", err)
					return
				}

				// URL de la API donde deseas hacer la solicitud POST
				apiURL := "http://" + connection + "/api/booking" // Ajusta la URL a tu entorno

				// Realiza la solicitud POST
				resp, err = http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
				if err != nil {
					fmt.Println("Error al hacer la solicitud POST:", err)
					return
				}

				defer resp.Body.Close()

				// Lee la respuesta del servidor si es necesario
				responseBody, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error al leer la respuesta del servidor:", err)
					return
				}
				//Unmarshall into FlightOrderResponse
				var response2 model.FlightOrderResponse
				err = json.Unmarshal(responseBody, &response2)
				if err != nil {
					fmt.Println("Error al leer la respuesta del servidor:", err)
					continue
				}

				//fmt.Println("Respuesta del servidor:", string(responseBody))
				fmt.Println("Reserva creada con éxito: ", response2.Data.ID)
			}

		case "2":
			{
				var id string
				fmt.Print("Ingrese el ID de la reserva: ")
				fmt.Scan(&id)

				url := fmt.Sprintf("http://"+connection+"/api/booking/%s", id)

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
				table.Append([]string{response.Data.FlightOffers[0].ID, response.Data.FlightOffers[0].Itineraries[0].Segments[0].Departure.At, response.Data.FlightOffers[0].Itineraries[0].Segments[0].Arrival.At, response.Data.FlightOffers[0].Itineraries[0].Segments[0].Aircraft.Code, response.Data.FlightOffers[0].Price.Total})
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
				for _, traveler := range response.Data.Travelers {
					passenger.Append([]string{traveler.Name.FirstName, traveler.Name.LastName})
				}
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
