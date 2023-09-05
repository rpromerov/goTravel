package model

type Search_Params struct {
	OriginLocationCode      string
	DestinationLocationCode string
	DepartureDate           string
	Adults                  string
	NonStop                 string
	CurrencyCode            string
	TravelClass             string
}
type FlightOfferPrice struct {
	Data []FlightOffer `json:"data"`
}

type FlightOfferResponse struct {
	Meta struct {
		Count int `json:"count"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"meta"`
	Data []FlightOffer `json:"data"`
}

type FlightOffer struct {
	Type                     string            `json:"type"`
	ID                       string            `json:"id"`
	Source                   string            `json:"source"`
	InstantTicketingRequired bool              `json:"instantTicketingRequired"`
	NonHomogeneous           bool              `json:"nonHomogeneous"`
	OneWay                   bool              `json:"oneWay"`
	LastTicketingDate        string            `json:"lastTicketingDate"`
	LastTicketingDateTime    string            `json:"lastTicketingDateTime"`
	NumberOfBookableSeats    int               `json:"numberOfBookableSeats"`
	Itineraries              []Itinerary       `json:"itineraries"`
	Price                    PriceDetails      `json:"price"`
	PricingOptions           PricingOptions    `json:"pricingOptions"`
	ValidatingAirlineCodes   []string          `json:"validatingAirlineCodes"`
	TravelerPricings         []TravelerPricing `json:"travelerPricings"`
}

type Itinerary struct {
	Duration string    `json:"duration"`
	Segments []Segment `json:"segments"`
}

type Segment struct {
	Departure     Location      `json:"departure"`
	Arrival       Location      `json:"arrival"`
	CarrierCode   string        `json:"carrierCode"`
	Number        string        `json:"number"`
	Aircraft      Aircraft      `json:"aircraft"`
	Operating     Operating     `json:"operating"`
	Duration      string        `json:"duration"`
	ID            string        `json:"id"`
	NumberOfStops int           `json:"numberOfStops"`
	Co2Emissions  []Co2Emission `json:"co2Emissions"`
}

type PriceDetails struct {
	Currency           string `json:"currency"`
	Total              string `json:"total"`
	Base               string `json:"base"`
	Fees               []Fee  `json:"fees"`
	GrandTotal         string `json:"grandTotal"`
	AdditionalServices []struct {
		Amount string `json:"amount"`
		Type   string `json:"type"`
	} `json:"additionalServices"`
}

type Fee struct {
	Amount string `json:"amount"`
	Type   string `json:"type"`
}

type PricingOptions struct {
	FareType                []string `json:"fareType"`
	IncludedCheckedBagsOnly bool     `json:"includedCheckedBagsOnly"`
}

type TravelerPricing struct {
	TravelerID           string       `json:"travelerId"`
	FareOption           string       `json:"fareOption"`
	TravelerType         string       `json:"travelerType"`
	Price                PriceDetails `json:"price"`
	FareDetailsBySegment []FareDetail `json:"fareDetailsBySegment"`
}

type FareDetail struct {
	SegmentID           string `json:"segmentId"`
	Cabin               string `json:"cabin"`
	FareBasis           string `json:"fareBasis"`
	BrandedFare         string `json:"brandedFare"`
	Class               string `json:"class"`
	IncludedCheckedBags struct {
		Quantity int `json:"quantity"`
	} `json:"includedCheckedBags"`
}

type FlightOffersPricing struct {
	Data         Data         `json:"data"`
	Dictionaries Dictionaries `json:"dictionaries"`
}

type Data struct {
	Type                string              `json:"type"`
	FlightOffers        []FlightOffer       `json:"flightOffers"`
	BookingRequirements BookingRequirements `json:"bookingRequirements"`
}

type Location struct {
	IATACode string `json:"iataCode"`
	At       string `json:"at"`
}

type Aircraft struct {
	Code string `json:"code"`
}

type Operating struct {
	CarrierCode string `json:"carrierCode"`
}

type Co2Emission struct {
	Weight     int    `json:"weight"`
	WeightUnit string `json:"weightUnit"`
	Cabin      string `json:"cabin"`
}

type Price struct {
	Currency        string `json:"currency"`
	Total           string `json:"total"`
	Base            string `json:"base"`
	Fees            []Fee  `json:"fees"`
	GrandTotal      string `json:"grandTotal"`
	BillingCurrency string `json:"billingCurrency"`
}

type FareDetailBySegment struct {
	SegmentID           string      `json:"segmentId"`
	Cabin               string      `json:"cabin"`
	FareBasis           string      `json:"fareBasis"`
	BrandedFare         string      `json:"brandedFare"`
	Class               string      `json:"class"`
	IncludedCheckedBags CheckedBags `json:"includedCheckedBags"`
}

type CheckedBags struct {
	Quantity int `json:"quantity"`
}

type BookingRequirements struct {
	EmailAddressRequired      bool `json:"emailAddressRequired"`
	MobilePhoneNumberRequired bool `json:"mobilePhoneNumberRequired"`
}

type Dictionaries struct {
	Locations map[string]LocationDetail `json:"locations"`
}

type LocationDetail struct {
	CityCode    string `json:"cityCode"`
	CountryCode string `json:"countryCode"`
}
