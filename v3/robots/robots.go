package robots

import (
	"bytes"
	"strconv"

	"github.com/Nigel2392/router/v3/request"
)

// Item in the robots.txt file.
type Listing struct {
	Allow      []string
	Disallow   []string
	UserAgent  string
	CrawlDelay int
}

// Options for generating the robots.txt file.
type Options struct {
	Rules   []*Listing
	SiteMap string
}

// Robots returns a handler that generates a robots.txt file.
func Robots(options *Options) func(r *request.Request) {
	return func(r *request.Request) {
		var buffer bytes.Buffer
		for i, listing := range options.Rules {
			if listing.UserAgent == "" {
				listing.UserAgent = "*"
			}
			buffer.WriteString("User-agent: " + listing.UserAgent + "\n")
			for _, allow := range listing.Allow {
				buffer.WriteString("Allow: " + allow + "\n")
			}
			for _, disallow := range listing.Disallow {
				buffer.WriteString("Disallow: " + disallow + "\n")
			}
			if listing.CrawlDelay > 0 {
				buffer.WriteString("Crawl-delay: " + strconv.Itoa(listing.CrawlDelay) + "\n")
			}
			if i < len(options.Rules)-1 {
				buffer.WriteString("\n")
			}
		}
		if options.SiteMap != "" {
			buffer.WriteString("\nSitemap: " + options.SiteMap + "\n")
		}
		r.Response.Header().Set("Content-Type", "text/plain")
		r.Response.Write(buffer.Bytes())
	}
}
