package fns

import "testing"

func TestRR(t *testing.T) {
	t.Run("value to MX", func(t *testing.T) {
		rr := RR{Value: "10 example.com"}
		mx := rr.ValueToMX()
		expectedMX := MXValue{
			Priority: 10,
			Domain:   "example.com",
		}

		if mx != expectedMX {
			t.Errorf("Test value to MX, expected mx %v, got: %v", expectedMX, mx)
		}
	})

	t.Run("value to SRV", func(t *testing.T) {
		rr := RR{Value: "10 10 8090 example.com"}
		srv := rr.ValueToSRV()
		expectedSRV := SRVValue{
			Priority: 10,
			Weight:   10,
			Port:     8090,
			Target:   "example.com",
		}

		if srv != expectedSRV {
			t.Errorf("Test value to SRV, expected srv %v, got: %v", expectedSRV, srv)
		}
	})

	t.Run("value to CAA", func(t *testing.T) {
		tests := []struct {
			rr          RR
			expectedCAA CAAValue
		}{
			{
				RR{Value: `0 issue "example.com"`},
				CAAValue{
					Flags: 0,
					Tag:   "issue",
					Value: `"example.com"`,
				},
			},
			{
				RR{Value: `0 iodef "mailto:security@example.com"`},
				CAAValue{
					Flags: 0,
					Tag:   "iodef",
					Value: `"mailto:security@example.com"`,
				},
			},
			{
				RR{Value: `0 issuewild "example.com; cansignhttpexchanges=yes"`},
				CAAValue{
					Flags: 0,
					Tag:   "issuewild",
					Value: `"example.com; cansignhttpexchanges=yes"`,
				},
			},
		}

		for i, test := range tests {
			caa := test.rr.ValueToCAA()
			if test.expectedCAA != caa {
				t.Errorf("Test %d, value to CAA, expected srv %v, got: %v", i, test.expectedCAA, caa)
			}
		}
	})
}
