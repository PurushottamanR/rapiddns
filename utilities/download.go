package utilities

import (
	"errors"
	"io"
	"net/http"
	"regexp"
)

func FetchRawData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error fetching raw data from RapidDNS")
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
