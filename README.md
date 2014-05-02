go_omtr - Adobe Omniture API client for Go
=======

Currently supports a limited set of the v1.4 Reporting API, documented at https://marketing.adobe.com/developer/en_US/documentation/analytics-reporting-1-4

## Usage

```go

import "github.com/ianmcmahon/go_omtr"


om_client := omtr.New("MyUsername:MyGroup", "shared-secret")

report_query := `{ "reportDescription": { "reportSuiteID":"my_report_suite", "metrics":[ { "id":"pageviews" } ] } }`

// queue a report:
report_id, err := om_client.QueueReport(report_query)

// reports take time to process, GetReport will return a Report Not Ready error until it's ready, at which point it will return the data:
data, err := om_client.GetReport(report_id)

// or, use the async Report() method, which takes a callback:
reportId, err := om_client.Report(report_query, func (data string) { fmt.Printf("Received data: %s\n", data) })
fmt.Printf("queued report: %d\n", reportId)
```



