package omtr

import (
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"strings"
	"encoding/json"
)



// returns status code, body as []byte, error
func (omcl *OmnitureClient) om_request(method, data string) (int, []byte, error) {
	endpoint := "https://api.omniture.com/admin/1.4/rest/?method=%s"

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf(endpoint, method), strings.NewReader(data))
	if err != nil { return -1, nil, err }

	req.Header.Add("X-WSSE", omcl.get_header())

	resp, err := client.Do(req)
	if err != nil { return -1, nil, err }

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil { return -1, nil, err }

	return resp.StatusCode, body, nil
}

func (omcl *OmnitureClient) QueueReport(query *ReportQuery) (int64, error) {
	bytes, err := json.Marshal(query)
	if err != nil { return -1, err }
	fmt.Printf("DEBUG: Marshalled report query to: %s\n", string(bytes))
	return omcl.QueueReportString(string(bytes))
}

// returns a report id which can be used at some point in the future to retrieve the report
func (omcl *OmnitureClient) QueueReportString(query string)  (int64, error) {
	status, b, err := omcl.om_request("Report.Queue", query)

	fmt.Printf("DEBUG: got back: %d, %s\n", status, b)

	response := queueReport_response{}

	err = json.Unmarshal(b, &response)
	if err != nil { return -1, err }

	fmt.Printf("DEBUG: unmarshaled: %v\n", response)

	return int64(response.ReportID), nil
}

func (omcl *OmnitureClient) GetReport(reportId int64)  ([]byte, error) {
	status, response, err := omcl.om_request("Report.Get", fmt.Sprintf("{ \"reportID\":%d }", reportId))
	if err != nil { return nil, err }

	// the api returns 400 if the report is not yet ready; in this case I'll parse the response as an error and return it
	if status == 400 {
		var ge getError
		err := json.Unmarshal(response, &ge)
		if err != nil { return nil, fmt.Errorf("Report.Get returned '%s'; error attempting to unmarshal to error structure: %v", string(response), err) }
		return nil, ge
	}

	fmt.Printf("DEBUG: got back: %d, %s\n", status, response)

	return response, err
}


/*
	Takes a report definition and a callback which will be called once the report has successfully been retrieved.
	Returns the reportId of the queued report or error
*/
func (omcl *OmnitureClient) Report(query *ReportQuery, success_callback func (string)) (int64, error) {
	rid, err := omcl.QueueReport(query)
	if err != nil { return -1, err }

	go omcl.wait_for_report_then_call(rid, success_callback)

	return rid, nil
}

func (omcl *OmnitureClient) ReportString(query string, success_callback func (string)) (int64, error) {
	rid, err := omcl.QueueReportString(query)
	if err != nil { return -1, err }

	go omcl.wait_for_report_then_call(rid, success_callback)

	return rid, nil
}

func (omcl *OmnitureClient) wait_for_report_then_call(rid int64, callback func(string)) {
	for {
		data, err := omcl.GetReport(rid)
		if err == nil {
			callback(string(data))
			return
		}
		time.Sleep(1 * time.Second)
	}
}
