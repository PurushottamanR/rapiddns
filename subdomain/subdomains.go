package subdomain

import (
)

type Options struct {
	Domain string
	All bool
	Pages int
	Verbose bool
}

var pattern = `(?m)<tr>\s*<th.*>(?P<row>.*)</th>\s*<td>(?P<subdomain>.*)</td>\s*<td>(<a\s*.*>)*\s*(?P<record>.*)\s*(</a>)*\s*</td>\s*<td>(?P<recordType>.*)</td>\s*<td>(?P<date>.*)</td>\s*</tr>`

func SubDomains(opts *Options) Records {

	src := NewSource("https://rapiddns.io/subdomain/%s?page=%d", "rapiddns", pattern)	
	return src.GetSubDomains(opts)
}
