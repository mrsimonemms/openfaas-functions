package function_test

import (
	"encoding/json"
	"flag"
	"handler/function"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	update = flag.Bool("update", false, "update the golden files of this test")
)

const feedXML = "testdata/feed.xml"

func TestParseXML(t *testing.T) {
	goldenFile := "testdata/feed.json"

	data, err := os.ReadFile(feedXML)
	assert.Nil(t, err)

	rates, err := function.ParseXML(data)
	assert.Nil(t, err)

	// Convert the feed to JSON format
	got, err := json.MarshalIndent(rates, "", "  ")
	require.NoError(t, err)

	if *update {
		err = os.WriteFile(goldenFile, []byte(got), 0600)
		require.NoError(t, err)
		return
	}

	content, err := os.ReadFile(goldenFile)
	assert.Nil(t, err)

	if diff := cmp.Diff(string(content), string(got)); diff != "" {
		t.Errorf("non-matching golden file (-want +got):\n%s", diff)
	}
}

func TestConvert(t *testing.T) {
	data, err := os.ReadFile(feedXML)
	require.NoError(t, err)

	rates, err := function.ParseXML(data)
	require.NoError(t, err)

	testCases := []struct {
		Name         string
		From         string
		To           string
		Amount       float64
		Date         *time.Time
		ExpectedRate float64
		Expected     float64
		Currencies   []function.Currency
	}{
		// {
		// 	Name:         "Same currency - whole number value",
		// 	From:         "EUR",
		// 	To:           "EUR",
		// 	Amount:       26,
		// 	Expected:     26,
		// 	ExpectedRate: 1,
		// },
		// {
		// 	Name:         "Same currency - decimal value",
		// 	From:         "EUR",
		// 	To:           "EUR",
		// 	Amount:       10.56,
		// 	Expected:     10.56,
		// 	ExpectedRate: 1,
		// },
		{
			Name:         "Based currency and non-base - whole number value",
			From:         "GBR",
			To:           "EUR",
			Amount:       25,
			Expected:     25,
			ExpectedRate: 1,
			Currencies:   rates.Currencies,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.Currencies == nil {
				testCase.Currencies = make([]function.Currency, 0)
			}

			c := function.ConversionRates{
				Currencies: testCase.Currencies,
			}

			r, err := c.Convert(testCase.From, testCase.To, testCase.Amount, testCase.Date)
			require.NoError(t, err)

			assert.Equal(t, r, &function.Conversion{
				Date: time.Now().UTC().Truncate(24 * time.Hour),
				From: function.CurrencyValue{
					Currency: testCase.From,
					Value:    testCase.Amount,
				},
				To: function.CurrencyValue{
					Currency: testCase.To,
					Value:    testCase.Expected,
				},
				Rate: testCase.Expected,
			})
		})
	}
}
