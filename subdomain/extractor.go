package subdomain

type Extractor interface {
	Extract([][]string) (Records, error)
}
