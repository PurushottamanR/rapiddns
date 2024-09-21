package subdomain


import (
	"fmt"
	"errors"
	"io"
	"net/http"
	"regexp"
	"time"
)

type Source struct {
	Name string
	Pattern string
}

func NewSource(name string, pattern string) *Source {
	
	src := Source{
		Name: name,
		Pattern: pattern,
	}
	
	return &src
}

func (s *Source) Fetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error fetching raw data from RapidDNS")
	}
	
	if resp.StatusCode == 429 {
	
		time.Sleep(time.Second * 15)
		resp, err = http.Get(url)
		if err != nil {
			return "", errors.New("Error fetching raw data from RapidDNS")
		}

	} 
	
	if resp.StatusCode != 200 {
		//after second attempt
		if resp.StatusCode == 429 {
			return "", errors.New(fmt.Sprintf("Response code: %d, try after sometime...", resp.StatusCode))
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

func (s *Source) GetSubDomains(url string) *SubDomainResults {
	results := &SubDomainResults {
		records: make(Records, 0, 100),
		Err: nil,
	}
	
	resp, err := s.Fetch(url)
	if err != nil {
		results.Err = err
		return results
	}
	
	matches := s.Match(resp, s.Pattern)
	records := s.Extract(matches)
	results.AddRecords(records)
	
	return results
}
