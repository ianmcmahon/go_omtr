package omtr

import (
	"encoding/json"
	"fmt"
	"time"
)

// takes a query object and returns a reportId which can be used to fetch the report in the future
func (omcl *OmnitureClient) QueueReport(query *ReportQuery) (int64, error) {
	bytes, err := json.Marshal(query)
	if err != nil {
		return -1, err
	}

	return omcl.QueueReportRaw(string(bytes))
}

// takes a query string (json) and returns a reportId which can be used to fetch the report in the future
func (omcl *OmnitureClient) QueueReportRaw(query string) (int64, error) {
	_, b, err := omcl.om_request("Report.Queue", query)
	if err != nil {
		return -1, err
	}

	response := queueReport_response{}

	err = json.Unmarshal(b, &response)
	if err != nil {
		return -1, err
	}

	return int64(response.ReportID), nil
}

// takes a reportId and returns a raw byteslice of json data, or error, including the Report Not Ready error.
func (omcl *OmnitureClient) GetReportRaw(reportId int64) ([]byte, error) {
	status, response, err := omcl.om_request("Report.Get", fmt.Sprintf("{ \"reportID\":%d }", reportId))
	if err != nil {
		return nil, err
	}

	// the api returns 400 if the report is not yet ready; in this case I'll parse the response as an error and return it
	if status == 400 {
		var ge getError
		err := json.Unmarshal(response, &ge)
		if err != nil {
			return nil, fmt.Errorf("Report.Get returned '%s'; error attempting to unmarshal to error structure: %v", string(response), err)
		}
		return nil, ge
	}

	return response, err
}

func (omcl *OmnitureClient) GetReport(reportId int64) (*ReportResponse, error) {
	bytes, err := omcl.GetReportRaw(reportId)
	if err != nil {
		return nil, err
	}

	resp := &ReportResponse{}

	err = json.Unmarshal(bytes, resp)

	resp.TimeRetrieved = time.Now()

	return resp, err
}

/*
	Takes a report definition and a callback which will be called once the report has successfully been retrieved.
	Returns the reportId of the queued report or error
*/
func (omcl *OmnitureClient) Report(query *ReportQuery, success_callback func(*ReportResponse, error)) (int64, error) {
	rid, err := omcl.QueueReport(query)
	if err != nil {
		return -1, err
	}

	go omcl.wait_for_report_then_call(rid, success_callback)

	return rid, nil
}

func (omcl *OmnitureClient) wait_for_report_then_call(rid int64, callback func(*ReportResponse, error)) {
	for {
		response, err := omcl.GetReport(rid)

		if err == nil {
			callback(response, nil)
			return
		}

		if _, ok := err.(getError); !ok { // type assert that err is a getError, and execute following block if it's not
			// getError is a "good" error; anything else isn't.
			callback(nil, err)
			return
		}

		time.Sleep(1 * time.Second)
	}
}
