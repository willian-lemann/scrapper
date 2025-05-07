package external

import (
	"log"
	"os"
	"scrapper/utils"

	"github.com/mendableai/firecrawl-go"
)

func Firecrawl() string {
	utils.LoadEnvs()

	app, err := firecrawl.NewFirecrawlApp(os.Getenv("FIRECRAWL_API_KEY"), "")
	if err != nil {
		log.Fatalf("Failed to initialize FirecrawlApp: %v", err)
	}

	scrapeResult, err := app.ScrapeURL("https://casaimoveisimb.com.br/imoveis/", nil)
	if err != nil {
		log.Fatalf("Failed to scrape URL: %v", err)
	}
	return scrapeResult.Markdown
}
