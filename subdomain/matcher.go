package subdomain

type Matcher interface {
	Match(httpResp string, pattern string) [][]string
}
