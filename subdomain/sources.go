package subdomain


import (
	"fmt"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

func (s *Source) GetTotalRecords() int {
	resp, _ := s.Fetch(fmt.Sprintf(s.GetURLFormat(), s.Opts.Domain, s.Opts.Pages))
	var page int = 1
	pagePattern := `class="page-link" href="/.*/[^"]+\?page=(?P<page>\d+)">`
	re := regexp.MustCompile(pagePattern)
	
	matches := re.FindAllStringSubmatch(resp, -1)
	if len(matches) > 0 {
		pageGroup := re.SubexpIndex("page")
		page, _ = strconv.Atoi(matches[0][pageGroup])
	}
	return page
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

func (s *Source) GetResultsFromPageThreaded(urlsCh chan string, resultsCh chan Records) {
	
	for url := range urlsCh {
	
		resp, err := s.Fetch(url)
		if err != nil && s.Opts.Verbose {
			log.Println(err)
			resultsCh <- Records{}
			continue
		}
		matches := s.Match(resp, s.Pattern)
		records := s.Extract(matches)
	
		resultsCh <- records
	}
}

func (s *Source) GetResultsThreaded(urls []string) Records {
	urlsCh := make(chan string, 20)
	resultsCh := make(chan Records)
	
	var records Records = make(Records, 0, 100)
	
	for i := 0; i < s.Opts.Threads; i++ {
		go s.GetResultsFromPageThreaded(urlsCh, resultsCh)
	}
	
	go func(){
		for _, url := range urls {
			urlsCh <- url
		}	
	}()
	
	
	for range urls {
		r := <- resultsCh
		records = append(records, r...)
		fmt.Printf("Retrieved total records: %d     \r", len(records))
	}
	
	close(urlsCh)
	close(resultsCh)
	
	return records
}


func (s *Source) GetResultsFromPage(url string) Records {
	resp, err := s.Fetch(url)
	if err != nil && s.Opts.Verbose {
		log.Println(err)
	}
	matches := s.Match(resp, s.Pattern)
	records := s.Extract(matches)
	
	return records
}

func (s *Source) GetResults(url string) Records {
	 return s.GetResultsFromPage(url)
}


func (s *Source) GetSubDomains() Records {
	var urls []string = []string{}
	var records Records
	
	if s.Opts.Pages > 1 {
		
		for page := 1; page <= s.Opts.Pages; page++ {
			u := fmt.Sprintf(s.GetURLFormat(), s.Opts.Domain, page)
			urls = append(urls, u)
		}
		
		records = s.GetResultsThreaded(urls)
		
		
	} else if s.Opts.All {
	
		pages := s.GetTotalRecords()
		for page := 1; page <= pages; page++ {
			u := fmt.Sprintf(s.GetURLFormat(), s.Opts.Domain, page)
			urls = append(urls, u)
		}
		
		records = s.GetResultsThreaded(urls)
		
		
	} else if s.Opts.Total {
	
		fmt.Println("Total pages:", s.GetTotalRecords())
		
	} else {
	
		url := fmt.Sprintf(s.GetURLFormat(), s.Opts.Domain, s.Opts.Pages) 
		
		records = s.GetResults(url)
		
		fmt.Printf("Retrieved total records: %d     \r", len(records))
	}
	
	return records
}
