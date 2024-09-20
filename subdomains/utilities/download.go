package utilities

import (
	"fmt"
	"errors"
	"io"
	"net/http"
	"regexp"
	"time"
)

func FetchRawResponse(url string) (*http.Response, error) {
	return http.Get(url)
}

func FetchRawData(url string) (string, error) {
	resp, err := FetchRawResponse(url)
	if err != nil {
		return "", errors.New("Error fetching raw data from RapidDNS")
	}
	
	if resp.StatusCode == 429 {
	
		time.Sleep(time.Second * 10)
		resp, err := FetchRawResponse(url)
		if err != nil {
			return "", errors.New("Error fetching raw data from RapidDNS")
		}
		//after second attempt
		if resp.StatusCode == 429 {
			return "", errors.New(fmt.Sprintf("Response code: %d, try after sometime...", resp.StatusCode))
		}
		
	} else if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Response code: %d", resp.StatusCode))
	}
	
	defer resp.Body.Close()
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading response")
	}
	return string(rawData), nil
}

func ExtractRecords(httpResp string, pattern string) [][]string {
	
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
