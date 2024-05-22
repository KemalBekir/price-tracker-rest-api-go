package scraper

// func ScrapeAmazon(url string) (model.Searches, error) {
// 	c := colly.NewCollector(
// 		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
// 	)

// 	var item model.Searches
// 	c.OnHTML("#title", func(e *colly.HTMLElement){
// 		item.ITEM_NAME = e.Text
// 	} )

// 	c.OnHTML(".a-price-whole",func(e *colly.HTMLElement){
// 		item.
// 	})
// }
