package subdomains

import (
	"fmt"
	
	"github.com/PurushottamanR/rapiddns/subdomains/utilities"
)

type Options struct {
	Domain string
	All bool
	Page int
	Verbose bool
}

type Record struct {
	ID string
	Subdomain string
	Value string
	RecType string
	Date string
}

type Records []Record

func (r Record) String() string {
	return fmt.Sprintf("%s %s %s %s", r.Subdomain, r.Value, r.RecType, r.Date)
}

func (R Records) String() string {
	var recs string = ""
	for _, record := range R {
		recs += record.String() + "\n"
	}
	return recs
}

var pattern = `(?m)<tr>\s*<th.*>(?P<row>.*)</th>\s*<td>(?P<subdomain>.*)</td>\s*<td>(<a\s*.*>)*\s*(?P<record>.*)\s*(</a>)*\s*</td>\s*<td>(?P<recordType>.*)</td>\s*<td>(?P<date>.*)</td>\s*</tr>`
var subdomainURL = "https://rapiddns.io/subdomain/%s?page=%d"


func FetchandExtract(url string) (Records, error) {
	var records Records = make(Records, 0, 100) 
	
	httpResp, err := utilities.FetchRawData(url)
	if err != nil {
		return records, err
	}

	matches := utilities.ExtractRecords(httpResp, pattern)
		
	
    	// Print the matched text content
    	for _, match := range matches {
    		record := Record{
    			ID: match[0], 
    			Subdomain: match[1], 
    			Value: match[2], 
    			RecType: match[3], 
    			Date: match[4],
    		}
    		records = append(records, record)
    	}
	return records, err
}

func SubDomains(opts *Options) (Records, error) {
	
	if !opts.All && opts.Page == 1 { 
	
		//Single and first page by default
		url := fmt.Sprintf(subdomainURL, opts.Domain, opts.Page)
		return FetchandExtract(url)
		 	
	}  else if !opts.All && opts.Page > 1 { 
		
		//multiple pages	
		var records Records = make(Records, 0, 100)
		var err error
		for page := 1; page <= opts.Page; page++ {
			url := fmt.Sprintf(subdomainURL, opts.Domain, page)
			pageRecords, err := FetchandExtract(url)
			if err != nil {
				return records, err
			}
			
			if len(pageRecords) > 0 {
				records = append(records, pageRecords...)
				fmt.Printf("[Page: %d][Nof Records: %d][Total Records: %d]\r", page, len(pageRecords), len(records))		
			}
		}
		return records, err
		
	} else {
		
		//every page - iterate until less than 100
		var records Records = make(Records, 0, 100)
		var err error
		for page := 1;; page++ {
			url := fmt.Sprintf(subdomainURL, opts.Domain, page)
			pageRecords, err := FetchandExtract(url)
			if err != nil {
				return records, err
			}
			
			if len(pageRecords) > 0 {
				records = append(records, pageRecords...)
				fmt.Printf("[Page: %d][Nof Records: %d][Total Records: %d]\r", page, len(pageRecords), len(records))		
			} else {
				break
			}
		}
		return records, err
	}
}
