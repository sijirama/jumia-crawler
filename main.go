package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"jumia-crawler/utils"
	"log"
	"time"
)

type Product struct {
	Name  string
	Price string
	URL   string
}

func scrapeWebsite(url string) ([]Product, error) {

	var products []Product
	headers := utils.GetRandomHeaders()

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.MaxDepth(0),
		colly.AllowedDomains("www.jumia.com.ng"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Second / 100,
	})

	c.OnRequest(func(r *colly.Request) {
		for k, v := range headers {
			r.Headers.Set(k, v)
		}
	})

	c.OnHTML("a.core[data-gtm-id][data-gtm-name]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Printf("Found Product link: %s\n", link)
		c.Visit(link)
	})

	c.OnHTML(".col10", func(e *colly.HTMLElement) {
		productName := e.ChildText("h1.-fs20")
		price := e.ChildText("span.-b.-ubpt.-tal.-fs24")
		productURL := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		if productName == "" || price == "" || productURL == "" {
			return
		}

		product := Product{
			Name:  productName,
			Price: price,
			URL:   productURL,
		}
		fmt.Printf("\nAppending: %s\n", product)
		products = append(products, product)
		fmt.Printf("We now have %d products\n", len(products))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping\n %s: %v\n\n", r.Request.URL, err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("Finished scraping %s\n", r.Request.URL)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	c.Wait()
	return products, nil
}

func main() {

	link := "https://www.jumia.com.ng/air-conditioners-d/haier-thermocool/"

	//"https://www.jumia.com.ng/haier-thermocool-1hp-genpal-inverter-air-conditioner-hsu-09lneb-03-white-3-years-warranty-132233135.html" // worked
	//"https://www.jumia.com.ng/mlp-appliances/"
	//"https://www.jumia.com.ng/"

	products, err := scrapeWebsite(link)

	if err != nil {
		log.Fatal("Error scraping website:", err)
	}

	fmt.Printf("Scraped %d products\n", len(products))

}
