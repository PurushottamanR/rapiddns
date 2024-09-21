package subdomain

type Fetcher interface {
	Fetch(url string) (string, error)
}
