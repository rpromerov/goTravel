package amadeus

import (
	"bytes"
	"encoding/json"
	"io"
	"model"
	"net/http"
)

func Flight_search(search_params model.Search_Params) model.FlightSearchResponse {
	println("Contacting Amadeus...")
	url := "https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=" + search_params.OriginLocationCode + "&destinationLocationCode=" + search_params.DestinationLocationCode + "&departureDate=" + search_params.DepartureDate + "&adults=" + search_params.Adults + "&nonStop=" + search_params.NonStop + "&currencyCode=" + search_params.CurrencyCode
	bearer := "Bearer " + "pS3shiu8IAHrjdzPuDkkxF5m8M1m"
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
	var response_body model.FlightSearchResponse
	err = json.Unmarshal([]byte(body), &response_body)
	if err != nil {
		println("Error " + resp.Status)
		panic(err)
	}

	return response_body
}
func Flight_pricing(req model.FlightData) model.FlightOfferResponse {
	url := "https://test.api.amadeus.com/v1/shopping/flight-offers/pricing"
	bearer := "Bearer " + "pS3shiu8IAHrjdzPuDkkxF5m8M1m"
	//unmarshall req to json
	bodyjson, err := json.Marshal(req)
	if err != nil {
		println("Error marshalling json " + err.Error())
		panic(err)
	}
	println(string(bodyjson))
	reqBody := bytes.NewBuffer([]byte(bodyjson))
	req2, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		panic(err)
	}

	req2.Header.Add("Authorization", bearer)
	req2.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req2)
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
	println(string(body))
	var response_body model.FlightOfferResponse
	err = json.Unmarshal([]byte(body), &response_body)
	if err != nil {
		println("Error " + resp.Status)
		panic(err)
	}

	return response_body
}
