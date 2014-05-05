go_omtr - Adobe Omniture API client for Go / Golang
=======

Currently supports a limited set of the v1.4 Reporting API, documented at https://marketing.adobe.com/developer/en_US/documentation/analytics-reporting-1-4

## Usage

```go

import "github.com/ianmcmahon/go_omtr"


om_client := omtr.New("MyUsername:MyGroup", "shared-secret")

// build a query either as a string:
report_query_string := `{ "reportDescription": { "reportSuiteID":"my_report_suite", "metrics":[ { "id":"pageviews" } ] } }`

// or as a query object:
report_query := &ReportQuery{
	ReportDesc: &ReportDescription {
		ReportSuiteID: "my_report_suite",
		Metrics: []*Metric{
			&Metric{"pageviews"},
		},
	},
}

// or use helper methods to construct a query:
report_query := omcl.Query("my_report_suite")
report_query.AddMetric("pageviews").AddMetric("event1").AddMetric("event2")
report_query.Granularity("hour")
report_query.DateFrom(time.Now().Add(7*24*time.Hour ))
report_query.DateTo(time.Now())

// queue your string query:
report_id, err := om_client.QueueReportString(report_query_string)

// or your object query:
report_id, err := om_client.QueueReport(report_query)

fmt.Printf("queued report: %d\n", report_id)

// reports take time to process, GetReport will return a Report Not Ready error until it's ready, at which point it will return the data:
data, err := om_client.GetReport(report_id)

// or, use the async Report() method, which takes a callback:
reportId, err := om_client.Report(report_query, func (data *ReportResponse, err error) { fmt.Printf("Received data: %v\n", data) })



```



