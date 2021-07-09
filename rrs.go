package fns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type rrsSDK struct {
	apiURL string
	client *http.Client
}

func NewRRsSDK(apiURL string) RRsRepository {
	return &rrsSDK{
		apiURL: apiURL,
		client: &http.Client{
			Timeout: time.Second * 3,
		},
	}
}

func (sdk rrsSDK) GetRRs(cond Condition) (rrs []RR) {
	req, err := http.NewRequest(http.MethodGet, sdk.apiURL, nil)
	if err != nil {
		log.Error("rrs-sdk new request error:", err)
		return
	}

	query := url.Values{}
	for param, value := range cond {
		query.Add(param, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = query.Encode()

	resp, err := sdk.client.Do(req)
	if err != nil {
		log.Error("rrs-sdk http client request error:", err)
		return
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("rrs-sdk read response body error:", err)
		return

	}

	type response struct {
		Code int64 `json:"code"`
		Data struct {
			List []RR `json:"list"`
		} `json:"data"`
	}

	var result response
	json.Unmarshal(b, &result)

	return result.Data.List
}
