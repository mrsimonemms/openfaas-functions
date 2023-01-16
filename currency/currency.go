package function

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BaseCurrency = "EUR" // The base currency is "1"
	FeedURL      = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
)

type ConversionRates struct {
	XMLName    xml.Name   `xml:"Envelope"`
	Subject    string     `xml:"subject"`
	Sender     Sender     `xml:"Sender"`
	Currencies []Currency `xml:"Cube>Cube"`
}

func (c *ConversionRates) findCurrency() {

}

// Creates the conversion struct with the target data
func (c *ConversionRates) Convert(from string, to string, amount float64, date *time.Time) (*Conversion, error) {
	now := time.Now().UTC()
	if date == nil {
		date = &now
	}

	// Build the conversion struct
	res := &Conversion{
		// Set the datetime to be start of day
		From: CurrencyValue{
			Currency: from,
			Value:    amount,
		},
		To: CurrencyValue{
			Currency: to,
		},
	}

	// Set the rates
	if from == to {
		// The currencies are the same - no conversion necessary
		res.Rate = 1
		// Date is unimportant, but stick to today for accuracy
		date = &now
	} else if from == BaseCurrency || to == BaseCurrency {
		// Different currencies, but one is base currency
		for _, k := range c.Currencies {
			fmt.Println(k.Date)
			fmt.Println(date)
		}
	}

	fmt.Println(from)
	fmt.Println(BaseCurrency)

	os.Exit(1)

	// Set to start of day
	res.Date = date.Truncate(24 * time.Hour)

	if err := res.Calculate(); err != nil {
		return nil, err
	}

	return res, nil
}

type Sender struct {
	Name string `xml:"name"`
}

type Currency struct {
	Date  CurrencyDate `xml:"time,attr"`
	Rates []Rate       `xml:"Cube"`
}

type CurrencyDate struct {
	time.Time
}

// Convert the string date into time.Time
func (c *CurrencyDate) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, err := time.Parse("2006-01-02", attr.Value)
	if err != nil {
		return err
	}
	c.Time = parse
	return nil
}

type Rate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

type Conversion struct {
	Date time.Time     `json:"date"`
	From CurrencyValue `json:"from"`
	To   CurrencyValue `json:"to"`
	Rate float64       `json:"rate"`
}

func (c *Conversion) Calculate() error {
	c.To.Value = c.From.Value * c.Rate
	return nil
}

type CurrencyValue struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

func GetFromFeed(client http.Client) (c *ConversionRates, err error) {
	resp, err := client.Get(FeedURL)
	if err != nil {
		return
	}

	defer func() {
		err = resp.Body.Close()
	}()
	if err != nil {
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return ParseXML(data)
}

func ParseXML(data []byte) (*ConversionRates, error) {
	var c *ConversionRates
	if err := xml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return c, nil
}
