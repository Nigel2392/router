# Example of robots.txt handler

var options = &robots.Options{
    Rules: []*robots.Listing{
        {
            Allow: []string{
                "/",
                "/about",
            },
            Disallow: []string{
                "/admin",
            },
            UserAgent: "*",
            CrawlDelay: 5,
        },
        {
            Disallow: []string{
                "/admin",
            },
            UserAgent: "Googlebot",
            CrawlDelay: 5,
        },
    }
    Sitemap: "https://example.com/sitemap.xml",
}

func main(){
    var r = router.NewRouter(false)
    var robotsHandler = robots.Robots(options)
    r.Get("/robots.txt", robotsHandler)
    r.Listen()
}
