package amadeus

import (
	"encoding/json"
	"io"
	"model"
	"net/http"
)

func Flight_search(search_params model.Search_Params) model.FlightOfferResponse {
	println("Contacting Amadeus...")
	url := "https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=" + search_params.OriginLocationCode + "&destinationLocationCode=" + search_params.DestinationLocationCode + "&departureDate=" + search_params.DepartureDate + "&adults=" + search_params.Adults + "&nonStop=" + search_params.NonStop + "&currencyCode=" + search_params.CurrencyCode
	bearer := "Bearer " + "JBrHMfCVad22Gq5BoJ2HqbQJqNt2"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("Error " + resp.Status)
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Error " + resp.Status)
		panic(err)
	}
	var response_body model.FlightOfferResponse
	err = json.Unmarshal([]byte(body), &response_body)
	if err != nil {
		println("Error " + resp.Status)
		panic(err)
	}

	return response_body
}
