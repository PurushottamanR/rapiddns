package subdomain

import (
	"fmt"
)

type Options struct {
	Domain string
	All bool
	Page int
	Verbose bool
}

var pattern = `(?m)<tr>\s*<th.*>(?P<row>.*)</th>\s*<td>(?P<subdomain>.*)</td>\s*<td>(<a\s*.*>)*\s*(?P<record>.*)\s*(</a>)*\s*</td>\s*<td>(?P<recordType>.*)</td>\s*<td>(?P<date>.*)</td>\s*</tr>`
var subdomainURL = "https://rapiddns.io/subdomain/%s?page=%d"

func SubDomains(opts *Options) *SubDomainResults {

	rapiddns := NewSource("rapiddns", pattern)	
	
	
	if !opts.All && opts.Page == 1 { 
	
		//Single and first page by default
		url := fmt.Sprintf(subdomainURL, opts.Domain, opts.Page)
		return rapiddns.GetSubDomains(url)

		 	
	}  else if !opts.All && opts.Page > 1 { 
		
		//multiple pages	
		results := &SubDomainResults{
			records: make(Records, 0, 100),
			Err: nil,
		}
		
		for page := 1; page <= opts.Page; page++ {
		
			url := fmt.Sprintf(subdomainURL, opts.Domain, page)
			
			res := rapiddns.GetSubDomains(url)
			if res.Err != nil {
				return res
			}
		
			if len(res.GetRecords()) > 0 {
				results.AddRecords(res.GetRecords())
				if opts.Verbose {
					fmt.Printf("[Page: %d][Nof Records: %d][Total Records: %d]\r", page, len(res.GetRecords()), len(results.GetRecords()))
				}
			} else {
				break
			}
		}
		return results
		
	} else {
		
		//every page - iterate until less than 100
		results := &SubDomainResults{
			records: make(Records, 0, 100),
			Err: nil,
		}
		
		for page := 1;; page++ {
		
			url := fmt.Sprintf(subdomainURL, opts.Domain, page)
			
			res := rapiddns.GetSubDomains(url)
			if res.Err != nil {
				return res
			}
		
			if len(res.GetRecords()) > 0 {
				results.AddRecords(res.GetRecords())
				if opts.Verbose {
					fmt.Printf("[Page: %d][Nof Records: %d][Total Records: %d]\r", page, len(res.GetRecords()), len(results.GetRecords()))
				}
			} else {
				break
			}
		}
		return results
	}
}
