package subdomain


import (
	"fmt"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Source struct {
	URLFormat string
	Name string
	Pattern string
	Opts *Options
}



func NewSource(urlformat string, name string, pattern string, opts *Options) *Source {
	
	src := Source{
		URLFormat: urlformat,
		Name: name,
		Pattern: pattern,
		Opts: opts,
	}
	
	return &src
}

func (s *Source) GetURLFormat() string {
	return s.URLFormat
}

func (s *Source) Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error fetching raw data from RapidDNS")
	}
	
	if resp.StatusCode == 429 {
	
		time.Sleep(time.Second * 30)
		resp, err = http.Get(url)
		if err != nil {
			return "", errors.New("Error fetching raw data from RapidDNS")
		}

	} 
	
	if resp.StatusCode != 200 {
		//after second attempt
		if resp.StatusCode == 429 {
			return "", errors.New(fmt.Sprintf("URL: %s, Response code: %d, try after sometime...", resp.Request.URL, resp.StatusCode))
		} else {
			return "", errors.New(fmt.Sprintf("Response code: %d", resp.StatusCode))
		}
	}
	
	defer resp.Body.Close()
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading response")
	}
	return string(rawData), nil
}

func (s *Source) Match(httpResp string, pattern string) [][]string {

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(httpResp, -1)
	
	records := make([][]string, len(matches))
	
	rowIndex := re.SubexpIndex("row")
	subdomainIndex := re.SubexpIndex("subdomain")
	recordIndex := re.SubexpIndex("record")
	recordTypeIndex := re.SubexpIndex("recordType")
	dateIndex := re.SubexpIndex("date")
	
    	// Print the matched text content
    	for i, match := range matches {
    		records[i] = append(records[i], match[rowIndex])
    		records[i] = append(records[i], match[subdomainIndex])
    		records[i] = append(records[i], match[recordIndex])
    		records[i] = append(records[i], match[recordTypeIndex])
    		records[i] = append(records[i], match[dateIndex])
    	}
    	return records

}

func (s *Source) Extract(matches [][]string) Records {
	var records Records = make(Records, 0, 100) 
		
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

func (s *Source) GetResultsFromPage(url string, recs chan Records) {
	resp, err := s.Fetch(url)
	if err != nil && s.Opts.Verbose {
		log.Println(err)
	}
	matches := s.Match(resp, s.Pattern)
	records := s.Extract(matches)
	
	recs <- records
}

func (s *Source) GetResults(urls []string) (chan Records) {
	recs := make(chan Records, len(urls)) 
	for _, url := range urls {
		go s.GetResultsFromPage(url, recs)
	}
	
	return recs
}


func (s *Source) GetSubDomains(opts *Options) Records {
	var urls []string = []string{}
	var records Records = make(Records, 0, 100)
	
	if opts.Pages > 1 {
		
		for page := 1; page <= opts.Pages; page++ {
			u := fmt.Sprintf(s.GetURLFormat(), opts.Domain, page)
			urls = append(urls, u)
		}
		
		recs := s.GetResults(urls)
		for range urls {
			r := <- recs
			records = append(records, r...)
			if opts.Verbose {
				fmt.Printf("Retrieved total records: %d     \r", len(records))
			}
		}
		
		
		close(recs)
		
		
	} else {
	
		urls = append(urls, fmt.Sprintf(s.GetURLFormat(), opts.Domain, opts.Pages)) 
		
		recs := s.GetResults(urls)
		records = <- recs
		close(recs)
		
		if opts.Verbose {
			fmt.Printf("Retrieved total records: %d     \r", len(records))
		}
	}
	
	return records
}
