package subdomain

import (
)

type Options struct {
	Domain string
	All bool
	Pages int
	Verbose bool
}

func NewOptions() *Options {
	opts := &Options {
		Domain: "",
		All: false,
		Pages: 1,
		Verbose: false,
	}
	
	return opts
}

var pattern = `(?m)<tr>\s*<th.*>(?P<row>.*)</th>\s*<td>(?P<subdomain>.*)</td>\s*<td>(<a\s*.*>)*\s*(?P<record>.*)\s*(</a>)*\s*</td>\s*<td>(?P<recordType>.*)</td>\s*<td>(?P<date>.*)</td>\s*</tr>`

func SubDomains(opts *Options) Records {

	src := NewSource("https://rapiddns.io/subdomain/%s?page=%d", "rapiddns", pattern, opts)	
	return src.GetSubDomains(opts)
}
