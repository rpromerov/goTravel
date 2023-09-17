package main

//hello world
import (
	"amadeus"
	"db_connector"
	"os"
	"time"

	"encoding/json"
	"fmt"
	"log"
	"model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func startApiTimer() {
	ticker := time.NewTicker(1500 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		amadeus.Get_api_key()
	}
}

func main() {
	fmt.Println("Getting environment variables...")
	err := godotenv.Load(".env")
	fmt.Println("Getting api keys...")
	amadeus.Get_api_key()
	//get a new key every 1500 secs
	go startApiTimer()
	if err != nil {
		println("Error getting environment variables " + err.Error())
	}
	fmt.Println("Connecting to MongoDB...")
	client, ctx, _, err := db_connector.Connect()
	if err != nil {
		println("Error connecting to MongoDB " + err.Error())
		panic(err)
	}
	_ = db_connector.Ping(client, ctx)

	fmt.Println("goTravel server starting...")
	app := fiber.New()
	app.Use(recover.New())
	app.Get("/api/search", func(c *fiber.Ctx) error {
		fmt.Println("GET /api/search")
		//get params from request
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
			return fiber.NewError(500, "{“message”: “Hubo un error al realizar la búsqueda”}")
		}
		return c.Send(json)
	})
	app.Post("/api/pricing", func(c *fiber.Ctx) error {
		fmt.Println("POST /api/pricing")
		// body as FlightOffersPricing
		var flightOfferPricing model.FlightData
		err := c.BodyParser(&flightOfferPricing)
		if err != nil {
			println("Error parsing body " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al retornar los precios”}")
		}
		//Fix missing carrierCode in segments
		for i := 0; i < len(flightOfferPricing.Data.FlightOffers); i++ {
			for j := 0; j < len(flightOfferPricing.Data.FlightOffers[i].Itineraries); j++ {
				for k := 0; k < len(flightOfferPricing.Data.FlightOffers[i].Itineraries[j].Segments); k++ {
					flightOfferPricing.Data.FlightOffers[i].Itineraries[j].Segments[k].Operating.CarrierCode = flightOfferPricing.Data.FlightOffers[i].Itineraries[j].Segments[k].CarrierCode
				}
			}
		}
		// call amadeus
		pricing_response := amadeus.Flight_pricing(flightOfferPricing)
		// marshall response to json
		json, err := json.Marshal(pricing_response)
		if err != nil {
			println("Error marshalling json " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al retornar los precios”}")
		}
		return c.Send(json)

	})
	app.Post("/api/booking", func(c *fiber.Ctx) error {
		fmt.Println("POST /api/booking")
		//body as FlightOrder
		var flightOrder model.FlightOrderResponse
		err := c.BodyParser(&flightOrder)
		if err != nil {
			println("Error parsing body " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al crear la reserva”}")
		}
		response := amadeus.Book_flight(flightOrder)
		order_json, err := json.Marshal(response)
		if err != nil {
			println("Error marshalling json " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al crear la reserva”}")
		}
		//store in mongodb
		_, err = db_connector.InsertOne(client, ctx, "reservations", response.Data)
		if err != nil {
			println("Error inserting reservation " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al crear la reserva”}")
		}
		return c.Send(order_json)
	})
	app.Get("/api/booking/:id", func(c *fiber.Ctx) error {
		fmt.Println("GET /api/booking/:id")
		//get id from request
		id := c.Params("id")
		fmt.Println(id)
		//get from mongodb
		filter := bson.M{"id": id}
		result, err := db_connector.GetOne(client, ctx, "reservations", filter)
		if err != nil {
			println("Error getting reservation " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al recuperar la reserva”}")
		}
		//mongo result as FlightOrderData
		var response model.FlightOrderData
		err = result.Decode(&response)
		if err != nil {
			println("Error decoding reservation " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al recuperar la reserva”}")
		}
		//wrap in FlightOrderResponse
		var response2 model.FlightOrderResponse
		response2.Data = response
		//marshall response to json
		json, err := json.Marshal(response2)
		if err != nil {
			println("Error marshalling json " + err.Error())
			return fiber.NewError(500, "{“message”: “Hubo un error al recuperar la reserva”}")
		}

		return c.Send(json)
	})
	log.Fatal(app.Listen(os.Getenv("SERVER") + ":" + os.Getenv("PORT")))
}
