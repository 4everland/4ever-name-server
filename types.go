package fns

import (
	"strconv"
	"strings"
)

type RRsRepository interface {
	GetRRs(cond Condition) []RR
}

type Condition map[string]interface{}

type RR struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	TTL    uint32 `json:"ttl"`
}

type SRVValue struct {
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   string
}

type CAAValue struct {
	Flags uint8
	Tag   string
	Value string
}

type MXValue struct {
	Priority uint16
	Domain   string
}

func (rr RR) ValueToSRV() (v SRVValue) {
	srvValue := strings.Split(rr.Value, " ")
	if len(srvValue) != 4 {
		return
	}

	var err error
	if v.Priority, err = rr.strToUint16(srvValue[0]); err != nil {
		return
	}

	if v.Weight, err = rr.strToUint16(srvValue[1]); err != nil {
		return
	}

	if v.Port, err = rr.strToUint16(srvValue[2]); err != nil {
		return
	}

	v.Target = srvValue[3]

	return
}

func (rr RR) ValueToCAA() (v CAAValue) {
	caaValue := strings.Split(rr.Value, " ")
	if len(caaValue) < 3 {
		return
	}

	priority, err := strconv.ParseUint(caaValue[0], 10, 8)
	if err != nil {
		return
	}
	v.Flags = uint8(priority)

	v.Tag = caaValue[1]
	v.Value = strings.Join(caaValue[2:], " ")

	return
}

func (rr RR) ValueToMX() (v MXValue) {
	var err error
	mxValue := strings.Split(rr.Value, " ")
	if len(mxValue) != 2 {
		return
	}

	if v.Priority, err = rr.strToUint16(mxValue[0]); err != nil {
		return
	}

	v.Domain = mxValue[1]

	return v
}

func (rr RR) strToUint16(str string) (uint16, error) {
	number, err := strconv.ParseUint(str, 10, 16)

	return uint16(number), err
}
