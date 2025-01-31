package iptools

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Options struct {
	IPaddr string
}

func NewOptions() *Options {

	opts := &Options{
		IPaddr: "",
	}

	return opts
}

type Datacenter struct {
	Datacenter         string `json:"datacenter"`
	Network            string `json:"network"`
	Country            string `json:"country"`
	Name               string `json:"name"`
	Region             string `json:"region"`
	City               string `json:"city"`
	RegionID           int    `json:"regionId"`
	Platform           string `json:"platform"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
	SystemService      string `json:"systemService"`
}

type Company struct {
	Name        string `json:"name"`
	AbuserScore string `json:"abuser_score"`
	Domain      string `json:"domain"`
	Type        string `json:"type"`
	Network     string `json:"network"`
	Whois       string `json:"whois"`
}

type Abuse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type Asn struct {
	Asn         int    `json:"asn"`
	AbuserScore string `json:"abuser_score"`
	Route       string `json:"route"`
	Descr       string `json:"descr"`
	Country     string `json:"country"`
	Active      bool   `json:"active"`
	Org         string `json:"org"`
	Domain      string `json:"domain"`
	Abuse       string `json:"abuse"`
	Type        string `json:"type"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	Rir         string `json:"rir"`
	Whois       string `json:"whois"`
}

type Location struct {
	Continent     string    `json:"continent"`
	Country       string    `json:"country"`
	CountryCode   string    `json:"country_code"`
	State         string    `json:"state"`
	City          string    `json:"city"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Zip           string    `json:"zip"`
	Timezone      string    `json:"timezone"`
	LocalTime     time.Time `json:"local_time"`
	LocalTimeUnix int       `json:"local_time_unix"`
	IsDst         bool      `json:"is_dst"`
}

type IPDetails struct {
	IP           string     `json:"ip"`
	Rir          string     `json:"rir"`
	IsBogon      bool       `json:"is_bogon"`
	IsMobile     bool       `json:"is_mobile"`
	IsCrawler    bool       `json:"is_crawler"`
	IsDatacenter bool       `json:"is_datacenter"`
	IsTor        bool       `json:"is_tor"`
	IsProxy      bool       `json:"is_proxy"`
	IsVpn        bool       `json:"is_vpn"`
	IsAbuser     bool       `json:"is_abuser"`
	Datacenter   Datacenter `json:"datacenter"`
	Company      Company    `json:"company"`
	Abuse        Abuse      `json:"abuse"`
	Asn          Asn        `json:"asn"`
	Location     Location   `json:"location"`
	ElapsedMs    float64    `json:"elapsed_ms"`
}

func (d *Datacenter) String() string {
	return fmt.Sprintf("Datacenter: %s\nNetwork: %s\nCountry: %s\nRegion: %s\nCity: %s\nPlatform: %s\nService: %s\nSystemService: %s\n", 
	d.Datacenter, 
	d.Network, 
	d.Country,
	d.Region,
	d.City,
	d.Platform,
	d.Service,
	d.SystemService)
}

func (c *Company) String() string {
	return fmt.Sprintf("Name: %s\nAbuserScore: %s\nDomain: %s\nType: %s\nNetwork: %s\n", c.Name, c.AbuserScore, c.Domain, c.Type, c.Network)
}

func (a *Abuse) String() string {
	return fmt.Sprintf("Name: %s\nAddress: %s\nEmail: %s\nPhone %s\n", a.Name, a.Address, a.Email, a.Phone)
}

func (a *Asn) String() string {
	return fmt.Sprintf("ASN: %d\nAbuserScore: %s\nRoute: %s\nDescription: %s\nCountry: %s\nActive: %t\nOrganisation: %s\nDomain: %s\nAbuse: %s\nType: %s\n", 
	a.Asn, 
	a.AbuserScore, 
	a.Route, 
	a.Descr, 
	a.Country, 
	a.Active, 
	a.Org, 
	a.Domain, 
	a.Abuse, 
	a.Type)
}

func (l *Location) String() string {
	return fmt.Sprintf("Continent: %s\nCountry: %s\nState: %s\nCity: %s\n", l.Continent, l.Country, l.State, l.City)
}

func (i *IPDetails) String() string {
	return fmt.Sprintf("\n%s\n%s\n%s\n%s\n", i.Datacenter.String(), i.Company.String(), i.Asn.String(), i.Location.String())
}

func GetIPDetails(options *Options) *IPDetails {
	url := fmt.Sprintf("https://<url>/?q=%s", options.IPaddr)
	var ipdetails *IPDetails = &IPDetails{}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(data, ipdetails)
	if err != nil {
		log.Fatalln(err)
	}

	return ipdetails
}
