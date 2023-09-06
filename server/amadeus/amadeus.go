package amadeus

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"model"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Flight_search(search_params model.Search_Params) model.FlightSearchResponse {
	url := "https://test.api.amadeus.com/v2/shopping/flight-offers?originLocationCode=" + search_params.OriginLocationCode + "&destinationLocationCode=" + search_params.DestinationLocationCode + "&departureDate=" + search_params.DepartureDate + "&adults=" + search_params.Adults + "&nonStop=" + search_params.NonStop + "&currencyCode=" + search_params.CurrencyCode
	bearer := "Bearer " + os.Getenv("api_key")
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
	bearer := "Bearer " + os.Getenv("api_key")
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
func Book_flight(order model.FlightOrderResponse) model.FlightOrderResponse {
	url := "https://test.api.amadeus.com/v1/booking/flight-orders"
	bearer := "Bearer " + os.Getenv("api_key")
	//unmarshall req to json
	bodyjson, err := json.Marshal(order)
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
	var response_body model.FlightOrderResponse
	err = json.Unmarshal([]byte(body), &response_body)
	if err != nil {
		println("Error " + resp.Status)
		panic(err)
	}

	return response_body
}
func Get_api_key() {
	urlPath := "https://test.api.amadeus.com/v1/security/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("client_secret", os.Getenv("SECRET_ID"))

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlPath, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	os.Setenv("api_key", strings.ReplaceAll(strings.Split(strings.Split(string(body), ",")[5], ":")[1], "\"", ""))

}
