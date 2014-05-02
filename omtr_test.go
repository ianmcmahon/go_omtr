package omtr

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func OmClient() *OmnitureClient {
	return New(os.Getenv("OM_USERNAME"), os.Getenv("OM_SECRET"))
}

func SampleQuery() *ReportQuery {
	return &ReportQuery{
		ReportDesc: &ReportDescription{
			ReportSuiteID:   "cnn-adbp-domestic",
			DateGranularity: "day",
			Date:            "2014-05-02",
			Metrics: []*Metric{
				&Metric{"pageviews"},
				&Metric{"event32"},
				&Metric{"event1"},
			},
		},
	}
}

func TestJsonNumberAsInt(t *testing.T) {
	var qr queueReport_response
	raw := `{"reportID": 42}`
	err := json.Unmarshal([]byte(raw), &qr)
	if err != nil {
		t.Error(err)
	}
}

func TestJsonNumberAsString(t *testing.T) {
	var qr queueReport_response
	raw := `{"reportID": "42"}`
	err := json.Unmarshal([]byte(raw), &qr)
	if err != nil {
		t.Error(err)
	}
}

func TestJsonGetError(t *testing.T) {
	var ge getError
	raw := `{"error":"report_not_ready","error_description":"Report not ready","error_uri":null}`
	err := json.Unmarshal([]byte(raw), &ge)
	if err != nil {
		t.Error(err)
	}
}

func TestQueueReport(t *testing.T) {
	omcl := OmClient()

	resp, err := omcl.QueueReport(SampleQuery())
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("response: %d\n", resp)
}

func TestReport(t *testing.T) {
	omcl := OmClient()

	c := make(chan *ReportResponse, 2)

	reportId, err := omcl.Report(SampleQuery(), func(response *ReportResponse, err error) {
		if err != nil {
			t.Error(err)
			c <- nil
		} else {
			c <- response
		}
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Submitted report, reportId is %d\n", reportId)

	go func() {
		time.Sleep(30 * time.Second)
		c <- nil
	}()

	for {
		data := <-c

		if data == nil {
			t.Errorf("Timed out waiting for report %d", reportId)
		}
		fmt.Printf("received data: %s\n", data)
		return
	}
}
