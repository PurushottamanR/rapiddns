package utilities

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var pattern = `<td[^>]*>(?:<a[^>]*>)?([^<]*)`

type domain struct {
	hostname string
}

func NewDomain(hostname string) *domain {
	
	return &domain{
		hostname: hostname,
	}
}

func (d *domain) GetSubDomains() string {
	page := 1
	url := fmt.Sprintf("https://rapiddns.io/subdomain/%s?page=%d", d.hostname, page)
	raw, err := fetchRawData(url)
	if err != nil {
		log.Fatalln(err)
	}
	
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(raw, -1)
	
	data := ""
    	// Print the matched text content
    	for _, match := range matches {
        	// Trim whitespace and print the capture group
        	record := ""
        	value := strings.TrimSpace(match[1])
        	record += value + " "
        	data += record + "\n"
    	}
	return data
}
