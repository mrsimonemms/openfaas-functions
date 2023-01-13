package function

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

const (
	FeedURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
)

type ConversionRates struct {
	XMLName    xml.Name   `xml:"Envelope"`
	Subject    string     `xml:"subject"`
	Sender     Sender     `xml:"Sender"`
	Currencies []Currency `xml:"Cube>Cube"`
}

type Sender struct {
	Name string `xml:"name"`
}

type Currency struct {
	Date  string `xml:"time,attr"`
	Rates []Rate `xml:"Cube"`
}

// Convert perform the conversion based upon today's rate
func (c *Currency) Convert(from string, to string, amount float64) (*Conversion, error) {
	return c.ConvertDate(from, to, amount, time.Now())
}

// ConvertDate performs the conversion based upon the given day's rate
func (c *Currency) ConvertDate(from string, to string, amount float64, date time.Time) (*Conversion, error) {
	return nil, nil
}

type Rate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

type Conversion struct {
	Date  time.Time `json:"date"`
	From  string    `json:"from"`
	To    string    `json:"to"`
	Rate  float64   `json:"rate"`
	Value float64   `json:"value"`
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

	c, err = ParseXML(data)

	return
}

func ParseXML(data []byte) (*ConversionRates, error) {
	var c *ConversionRates
	if err := xml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return c, nil
}
