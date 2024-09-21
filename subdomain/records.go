package subdomain

import (
	"fmt"
)

type Record struct {
	ID string
	Subdomain string
	Value string
	RecType string
	Date string
}

type Records []Record

type SubDomainResults struct {
	records	   Records
	Err	   error
}

func (r *Record) String() string {
	return fmt.Sprintf("%s %s %s %s", r.Subdomain, r.Value, r.RecType, r.Date)
}

func (r Records) String() string {
	var recs string = ""
	for _, record := range r {
		recs += record.String() + "\n"
	}
	return recs
}

func (r *SubDomainResults) GetRecords() Records {
	return r.records
}

func (r *SubDomainResults) AddRecords(rec Records) {
	r.records = append(r.records, rec...)
}

func (r *SubDomainResults) String() string {
	if r.Err != nil {
		return r.Err.Error()
	} else {
		return r.records.String()
	}
}
