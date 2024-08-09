package subdomains

import (
	"fmt"
	"log"
	"time"
	
	"github.com/PurushottamanR/rapiddns/utilities"
)

type Record struct {
	ID string
	Subdomain string
	Value string
	RecType string
	Date string
}

type domain struct {
	hostname string
}

func (r Record) String() string {
	return fmt.Sprintf("%s %s %s %s", r.Subdomain, r.Value, r.RecType, r.Date)
}

var pattern = `(?m)<tr>\s*<th.*>(?P<row>.*)</th>\s*<td>(?P<subdomain>.*)</td>\s*<td>(<a\s*.*>)*\s*(?P<record>.*)\s*(</a>)*\s*</td>\s*<td>(?P<recordType>.*)</td>\s*<td>(?P<date>.*)</td>\s*</tr>`
var subdomainURL = "https://rapiddns.io/subdomain/%s?page=%d"

func NewDomain(hostname string) *domain {
	
	return &domain{
		hostname: hostname,
	}
}

func FetchandExtract(url string) []Record {
	var records []Record = make([]Record, 0, 100)
	httpResp, err := utilities.FetchRawData(url)
	if err != nil {
		log.Fatalln(err)
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
	return records
}

func (d *domain) SubDomains(all bool, page int) []Record {
	
	if !all {
		url := fmt.Sprintf(subdomainURL, d.hostname, page)
		return FetchandExtract(url)
	    	
	}  else {
		var records []Record = make([]Record, 0, 100)
		for page := 1;; page++ {
			url := fmt.Sprintf(subdomainURL, d.hostname, page)
			pageRecords := FetchandExtract(url)
			if len(pageRecords) > 0 {
				records = append(records, pageRecords...)
				fmt.Printf("Nof Records: %d\r", len(records))		
			} else {
				break
			}
			time.Sleep(time.Second * 1)		
		}
		return records	
	}
}
