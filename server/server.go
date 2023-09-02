package main

import (
	"fmt"
	"net/http"
)

type Flight struct {
	Type            string `json:"type"`
	Origin          string `json:"origin"`
	Destination     string `json:"destination"`
	DepartureDate   string `json:"departureDate"`
	ReturnDate      string `json:"returnDate"`
	Price           struct {
		Total string `json:"total"`
	} `json:"price"`
	Links struct {
		FlightDates   string `json:"flightDates"`
		FlightOffers  string `json:"flightOffers"`
	} `json:"links"`
}
type Response struct {
	Data []FlightDestination `json:"data"`
}


func main() {
	// Definir una ruta y asociarla con una función de controlador

	// ejemplo
	jsonData := `
		un json
	`
	var response Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		fmt.Println("Error al analizar JSON:", err)
		return
	}
	flights = response.Data

	// Definir rutas de la API.
	http.HandleFunc("/flights", getFlights)
	http.HandleFunc("/flight/", getFlightByDestination)

	// Iniciar el servidor en el puerto 8080.
	fmt.Println("Servidor API de vuelos en ejecución en el puerto 8080...")
	http.ListenAndServe(":8080", nil)

	// Definir rutas y controladores para los nuevos endpoints
    r.HandleFunc("/api/search", searchHandler).Methods("GET")
    r.HandleFunc("/api/pricing", pricingHandler).Methods("POST")
    r.HandleFunc("/api/booking", bookingHandler).Methods("POST")
    r.HandleFunc("/api/booking/{id}", getBookingHandler).Methods("GET")

    http.Handle("/", r)

    // Iniciar el servidor en el puerto 5000
    fmt.Println("Servidor en ejecución en el puerto 5000...")
    http.ListenAndServe(":5000", nil)
}

// Implementar los controladores para los nuevos endpoints
func searchHandler(w http.ResponseWriter, r *http.Request) {
    // Lógica de búsqueda de vuelos con Amadeus
    // ...
}

func pricingHandler(w http.ResponseWriter, r *http.Request) {
    // Lógica de obtención de precios de vuelos con Amadeus
    // ...
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {
    // Lógica de creación de reserva de vuelo con Amadeus
    // ...
}

func getBookingHandler(w http.ResponseWriter, r *http.Request) {
    // Lógica de obtención de reserva de vuelo con Amadeus
    // ...
}

func getFlights(w http.ResponseWriter, r *http.Request) {
	// Devuelve la lista de vuelos como respuesta JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flights)
}

func getFlightByDestination(w http.ResponseWriter, r *http.Request) {
	// Obtiene la destinación de la URL.
	destination := r.URL.Path[len("/flight/"):]

	// Busca vuelos por destinación en la lista.
	var matchingFlights []FlightDestination
	for _, flight := range flights {
		if flight.Destination == destination {
			matchingFlights = append(matchingFlights, flight)
		}
	}

	if len(matchingFlights) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matchingFlights)
	} else {
		// Si no se encuentra ningún vuelo para la destinación, devolver un error 404.
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No se encontraron vuelos para la destinación %s", destination)
	}
}
