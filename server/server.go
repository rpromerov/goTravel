package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"os"
	//"github.com/gin-gonic/gin"
	"github.com/olekukonko/tablewriter"
	"net/url"
	"strings"
	"io/ioutil"
)
type AccessTokenResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
}

type RespuestaJSON struct {
	Data []Flight `json:"data"`
}

// Estructura para representar la oferta de vuelo
type Flight struct {
	ID           	string `json:"id"`
	CarrierCode		string `json:"carrierCode"`
	Number 			string `json:"number"`	
	Departure       string `json:"departure"`
	Arrival		    string `json:"arrival"`
	Code			string `json:"code"`
	Total           string `json:"total"`
}

func obtenerNuevoToken() (string, error) {
    clientID := "rVLNX6dP527lBIpfBEyWGpKt92xhxjlz"
    clientSecret := "vOlBYNKISDUGamcp"

    if clientID == "" || clientSecret == "" {
        return "", fmt.Errorf("No se han configurado las credenciales de Amadeus.")
    }

    authURL := "https://test.api.amadeus.com/v1/security/oauth2/token"

    // Parámetros del formulario
    form := url.Values{}
    form.Add("grant_type", "client_credentials")
    form.Add("client_id", clientID)
    form.Add("client_secret", clientSecret)

    resp, err := http.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Error al obtener el token. Código de estado: %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var tokenResponse AccessTokenResponse
    if err := json.Unmarshal(body, &tokenResponse); err != nil {
        return "", err
    }

    return tokenResponse.AccessToken, nil
}

func main() {
	// Si opción es 1
	var origen, destino, fechaSalida, adultos string

	// Usuario ingresa datos
	fmt.Print("Aeropuerto de origen: ")
	fmt.Scan(&origen)

	fmt.Print("Aeropuerto de destino: ")
    fmt.Scan(&destino)

    fmt.Print("Fecha de salida (YYYY-MM-DD): ")
    fmt.Scan(&fechaSalida)

    fmt.Print("Cantidad de adultos: ")
    fmt.Scan(&adultos)

	// URL del endpoint de ofertas de vuelos
    url := "https://test.api.amadeus.com/v2/shopping/flight-offers"
    
    // Crear una solicitud HTTP GET con los parámetros
    req, err := http.NewRequest("GET", url, nil)
    
    // Agregar variables a la URL
    q := req.URL.Query()
    q.Add("originLocationCode", "ARI")
    q.Add("destinationLocationCode", "SCL")
    q.Add("departureDate", "2023-12-02")
    q.Add("adults", "1")
	// Agregar constantes
    q.Add("includedAirlineCodes", "H2,LA,JA")
    q.Add("nonStop", "true")
    q.Add("currencyCode", "CLP")
    q.Add("travelClass", "ECONOMY")
    req.URL.RawQuery = q.Encode()
    
    // Agregar el encabezado de autorización
	accessToken, err := obtenerNuevoToken()
    if err != nil {
        fmt.Println("Error al obtener el token:", err)
        return
    }
    req.Header.Set("Authorization", "Bearer " + accessToken)
    
    // Realizar la solicitud HTTP
    client := &http.Client{}
    resp2, err := client.Do(req)
    defer resp2.Body.Close()

    // Leer y mostrar la respuesta
    body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println("Error al deserializar JSON:", err)
		return
	}

	respuestaJSON := body2
	//fmt.Println(string(body2))

	// Deserializar la respuesta JSON en la estructura RespuestaJSON
	var respuesta RespuestaJSON
	Jerr := json.Unmarshal([]byte(respuestaJSON), &respuesta)
	if Jerr != nil {
		fmt.Println("Error al deserializar JSON:", err)
		return
	}

	// Crear una nueva tabla
	table := tablewriter.NewWriter(os.Stdout)
	// Encabezados de la tabla
	headers := []string{"VUELO", "NÚMERO", "HORA DE SALIDA", "HORA DE LLEGADA", "AVIÓN", "PRECIO TOTAL"}
	// Setear los encabezados de la tabla
	table.SetHeader(headers)


	// Acceder a los datos de las ofertas de vuelos
	for _, oferta := range respuesta.Data {
		// Acceder a los itinerarios de la oferta de vuelo
		row := []string{
			oferta.ID,
			oferta.CarrierCode + oferta.Number,
			oferta.Departure,
			oferta.Arrival,
			"A" + oferta.Code,
			oferta.Total,
		}
		table.Append(row)
	}
	// Imprimir la tabla en la salida estándar
	table.Render()
}
