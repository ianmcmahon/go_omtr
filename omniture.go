package omtr

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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
