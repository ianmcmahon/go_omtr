package omtr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	date_format	string = "2006-01-02"
)

// construct a Report Query for a reportsuite
func Query(reportSuite string) *ReportQuery {
	return &ReportQuery{
		&Description{
			ReportSuiteID: reportSuite,
		},
	}
}

// this way clients can either call omtr.Query (above) or omcl.Query
func (omcl *OmnitureClient) Query(reportSuite string) *ReportQuery {
	return Query(reportSuite)
}

func (q *ReportQuery) AddMetric(m string) *ReportQuery {
	q.ReportDescription.Metrics = append(q.ReportDescription.Metrics, &Metric{m})
	return q
}

func (q *ReportQuery) Granularity(g string) *ReportQuery {
	q.ReportDescription.DateGranularity = g
	return q
}

func (q *ReportQuery) Date(d time.Time) *ReportQuery {
	q.ReportDescription.Date = d.Format(date_format)
	return q
}

func (q *ReportQuery) DateFrom(d time.Time) *ReportQuery {
	q.ReportDescription.DateFrom = d.Format(date_format)
	return q
}

func (q *ReportQuery) DateTo(d time.Time) *ReportQuery {
	q.ReportDescription.DateTo = d.Format(date_format)
	return q
}

// returns status code, body as []byte, error
func (omcl *OmnitureClient) om_request(method, data string) (int, []byte, error) {
	endpoint := "https://api.omniture.com/admin/1.4/rest/?method=%s"

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf(endpoint, method), strings.NewReader(data))
	if err != nil {
		return -1, nil, err
	}

	req.Header.Add("X-WSSE", omcl.get_header())

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}

	return resp.StatusCode, body, nil
}
