package main

//hello world
import (
	"amadeus"
	"encoding/json"
	"fmt"
	"log"
	"model"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("goTravel server starting...")
	app := fiber.New()
	app.Get("/api/search", func(c *fiber.Ctx) error {
		//get params from request
		println(c.Query("originLocationCode"))
		originLocationCode := c.Query("originLocationCode")
		destinationLocationCode := c.Query("destinationLocationCode")
		departureDate := c.Query("departureDate")
		adults := c.Query("adults")
		nonStop := "true"
		currencyCode := "CLP"
		travelClass := "ECONOMY"

		//fill search struct
		search := model.Search_Params{
			OriginLocationCode:      originLocationCode,
			DestinationLocationCode: destinationLocationCode,
			DepartureDate:           departureDate,
			Adults:                  adults,
			NonStop:                 nonStop,
			CurrencyCode:            currencyCode,
			TravelClass:             travelClass,
		}
		search_response := amadeus.Flight_search(search)

		//marshall response to json
		json, err := json.Marshal(search_response)
		if err != nil {
			println("Error marshalling json " + err.Error())
			panic(err)
		}
		return c.Send(json)
	})
	app.Post("/api/pricing", func(c *fiber.Ctx) error {
		// body as FlightOffersPricing
		var flightOfferPricing model.FlightData
		err := c.BodyParser(&flightOfferPricing)
		if err != nil {
			println("Error parsing body " + err.Error())
			panic(err)
		}
		// call amadeus
		pricing_response := amadeus.Flight_pricing(flightOfferPricing)
		// marshall response to json
		json, err := json.Marshal(pricing_response)
		if err != nil {
			println("Error marshalling json " + err.Error())
			panic(err)
		}
		return c.Send(json)

	})

	log.Fatal(app.Listen(":5001"))
}
