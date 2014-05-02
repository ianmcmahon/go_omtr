package omtr

import (
	"os"
	"time"
	"fmt"
	"testing"
	"encoding/json"
)

func OmClient() *OmnitureClient {
	return New(os.Getenv("OM_USERNAME"), os.Getenv("OM_SECRET"))
}

func TestJsonNumberAsInt(t *testing.T) {
	var qr queueReport_response
	raw := `{"reportID": 42}`
	err := json.Unmarshal([]byte(raw), &qr)
	if err != nil { t.Error(err) }
}

func TestJsonNumberAsString(t *testing.T) {
	var qr queueReport_response
	raw := `{"reportID": "42"}`
	err := json.Unmarshal([]byte(raw), &qr)
	if err != nil { t.Error(err) }
}

func TestJsonGetError(t *testing.T) {
	var ge getError
	raw := `{"error":"report_not_ready","error_description":"Report not ready","error_uri":null}`
	err := json.Unmarshal([]byte(raw), &ge)
	if err != nil { t.Error(err) }
}

func TestQueueReport(t *testing.T) {
	omcl := OmClient()

	test_rept := "{\"reportDescription\":{\"reportSuiteID\":\"cnn-adbp-domestic\", \"dateFrom\":\"2014-04-20\", \"dateTo\":\"2014-04-21\", \"dateGranularity\":\"hour\", \"metrics\":[{\"id\":\"event32\"},{\"id\":\"event1\"},{\"id\":\"visitorsHourly\"},{\"id\":\"pageviews\"} ], \"currentData\":\"true\"} }"

	resp, err := omcl.QueueReport(test_rept)
	if err != nil { t.Error(err) }

	fmt.Printf("response: %d\n", resp)
}

func TestReport(t *testing.T) {
	omcl := OmClient()

	test_rept := "{\"reportDescription\":{\"reportSuiteID\":\"cnn-adbp-domestic\", \"dateFrom\":\"2014-04-20\", \"dateTo\":\"2014-04-21\", \"dateGranularity\":\"hour\", \"metrics\":[{\"id\":\"event32\"},{\"id\":\"event1\"},{\"id\":\"visitorsHourly\"},{\"id\":\"pageviews\"} ], \"currentData\":\"true\"} }"

	c := make(chan string, 2)

	reportId, err := omcl.Report(test_rept, func (data string) { c <- data })
	if err != nil { t.Error(err) }
	fmt.Printf("Submitted report, reportId is %d\n", reportId)

	go func () { 
		time.Sleep(30 * time.Second)
		c <- "error"
	}()

	for {
		data := <- c

		if data == "error" { t.Errorf("Timed out waiting for report %d", reportId) }
		fmt.Printf("received data: %s\n", data)
		return
	}
}