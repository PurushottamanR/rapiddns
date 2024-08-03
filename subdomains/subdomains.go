package subdomains

import (
	"fmt"
	"log"
	
	"github.com/PurushottamanR/rapiddns/utilities"
)

var pattern = `(?m)<tr>\s*<th.*>(?P<row>.*)</th>\s*<td>(?P<subdomain>.*)</td>\s*<td>(<a\s*.*>)*\s*(?P<record>.*)\s*(</a>)*\s*</td>\s*<td>(?P<recordType>.*)</td>\s*<td>(?P<date>.*)</td>\s*</tr>`
var subdomainURL = "https://rapiddns.io/subdomain/%s?page=%d"

type domain struct {
	hostname string
}

func NewDomain(hostname string) *domain {
	
	return &domain{
		hostname: hostname,
	}
}

func (d *domain) GetSubDomains() {
	page := 1
	url := fmt.Sprintf(subdomainURL, d.hostname, page)
	httpResp, err := utilities.FetchRawData(url)
	if err != nil {
		log.Fatalln(err)
	}

	records := utilities.ExtractRecords(httpResp, pattern)
	
    	// Print the matched text content
    	for _, record := range records {
    		fmt.Println(record[0], record[1], record[2], record[3], record[4])
    	}
}
