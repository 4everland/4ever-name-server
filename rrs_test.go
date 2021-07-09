package fns

import (
	"encoding/json"
	"github.com/miekg/dns"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRrsSDK(t *testing.T) {
	expectedRRs := []RR{
		{
			Name:   "@",
			Value:  "1.2.3.4",
			Domain: "example.com",
			TTL:    600,
		},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)

		resp := map[string]interface{}{
			"data": map[string][]RR{
				"list": expectedRRs,
			},
		}

		b, _ := json.Marshal(resp)
		res.Write(b)
	}))

	defer testServer.Close()

	repo := NewRRsSDK(testServer.URL)

	rrs := repo.GetRRs(Condition{
		"name": "example.com.",
		"type": dns.TypeA,
	})

	if !reflect.DeepEqual(expectedRRs, rrs) {
		t.Errorf("Test RRs Expected: '%v', got: '%v'", expectedRRs, rrs)
	}
}
