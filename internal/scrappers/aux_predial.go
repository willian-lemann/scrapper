package scrappers

import (
	"fmt"
	"log"
	"net/url"
	"scrapper/internal/repositories"
	"scrapper/internal/structs"
	"scrapper/utils"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func checkListing(listing structs.ListingItem) bool {
	if strings.Trim(listing.Type, " ") == "" {
		return false
	}
	if strings.Contains(listing.Type, "Galpão") ||
		strings.Contains(listing.Type, "Mônaco") ||
		strings.Contains(listing.Type, "Sobrado") ||
		strings.Contains(listing.Type, "Terreno") {
		return false
	}

	return true
}

func getNavLinks() []string {
	mainLink := "https://www.auxiliadorapredial.com.br/comprar/residencial/sc+imbituba?"

	var links []string

	quantityOfPages := 19
	for i := 1; i <= quantityOfPages; i++ {
		newLink := fmt.Sprintf("%spage=%d", mainLink, i)
		links = append(links, newLink)
	}

	return links
}

func initScrape(visitLink string) []structs.ListingItem {
	c := colly.NewCollector()

	listings := []structs.ListingItem{}

	c.OnHTML("#contentSide > div.MuiBox-root.css-164mavz > div > div.sc-a7b1d3df-0.cwVfwz", func(h *colly.HTMLElement) {
		h.ForEach("div div.content", func(_ int, h *colly.HTMLElement) {
			ref := strings.Replace(h.ChildText("div.footer div.ref span"), "ref: ", "", -1)
			link := fmt.Sprintf("https://www.auxiliadorapredial.com.br/imovel/venda/%s", ref)
			id, _ := strconv.Atoi(ref)
			address := h.ChildText("div.Location span")
			priceValue := strings.Trim(h.ChildText("div.content div.headContent div.total div div.oldValue"), " ")
			price := strings.Replace(priceValue, "R$", "", -1)
			listingType := h.ChildText("div.content div.headContent h4")
			bedrooms := 0
			bathrooms := 0
			parking := 0
			area := ""
			forSale := true

			h.ForEach("div.Details", func(i int, h *colly.HTMLElement) {
				h.ForEach("div", func(i int, h *colly.HTMLElement) {
					alt := h.ChildAttr("img", "alt")
					value := h.ChildText("span")
					if alt == "Metragem" {
						area = value
					}
					if alt == "Banheiros" {
						bathrooms, _ = strconv.Atoi(value)
					}
					if alt == "Garagens" {
						parking, _ = strconv.Atoi(value)
					}
					if alt == "Quartos" {
						bedrooms, _ = strconv.Atoi(value)
					}
				})
			})

			if strings.Contains(listingType, "Apartamento") {
				listingType = "Apartamento"
			}
			if strings.Contains(listingType, "Casa") {
				listingType = "Casa"
			}
			if strings.Contains(listingType, "Terreno") {
				listingType = "Terreno"
			}
			if strings.Contains(listingType, "Sobrado") {
				listingType = "Sobrado"
			}

			listing := structs.ListingItem{
				Link:      link,
				Id:        id,
				Image:     "",
				Address:   address,
				Price:     price,
				Type:      listingType,
				Bedrooms:  bedrooms,
				Bathrooms: bathrooms,
				Parking:   parking,
				ForSale:   forSale,
				Ref:       ref,
				Area:      area,
				Agency:    "aux_predial",
			}

			if checkListing(listing) {
				listings = append(listings, listing)
			}
		})
	})

	c.Visit(visitLink)
	return listings
}

func extractImageURL(value string) string {
	decodedURL, err := url.QueryUnescape(value)
	if err != nil {
		return ""
	}

	imageURL := strings.TrimPrefix(decodedURL, "/imovel/_next/image?url=")
	imageURL = strings.Split(imageURL, "&w=")[0]

	return imageURL
}

func insidePageScrape(listing structs.ListingItem) structs.ListingItem {
	c := colly.NewCollector()

	c.OnHTML("section.section-sobre-detalhe div#descricao", func(h *colly.HTMLElement) {
		id, _ := strconv.Atoi(utils.GetIDFromLink(h.Response.Request.URL.String()))
		if id == listing.Id {
			bodyHtml, err := h.DOM.Html()
			if err != nil {
				log.Fatal(err)
			}
			listing.Content = bodyHtml
		}
	})
	c.OnHTML("main#detalhe div.layout-control-detalhe section.exibicao-container div.guia div.exibicao-fotos-5 a", func(h *colly.HTMLElement) {
		url := h.Response.Request.URL.String()
		id, _ := strconv.Atoi(utils.GetIDFromLink(url))
		srcImage := h.ChildAttr("img", "src")

		if id == listing.Id {
			imagePhoto := extractImageURL(srcImage)

			if utils.CheckURLImage(imagePhoto) {
				listing.Photos = append(listing.Photos, structs.Photos{
					Href:          imagePhoto,
					ListingItemId: listing.Id,
				})
			}
		}
	})

	c.Visit(listing.Link)
	return listing
}

func scrapeWorker(link string, ch chan structs.ListingItem, w *sync.WaitGroup) {
	defer w.Done()

	listings := initScrape(link)

	for _, listing := range listings {
		newListing := insidePageScrape(listing)
		if len(newListing.Photos) == 0 {
			continue
		}
		newListing.Image = newListing.Photos[0].Href
		ch <- newListing
	}
}

func ExecuteAuxPredial(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Starting scrapping aux_predial")

	links := getNavLinks()

	resultch := make(chan structs.ListingItem)

	var w sync.WaitGroup

	for _, link := range links {
		w.Add(1)
		go scrapeWorker(link, resultch, &w)
	}

	go func() {
		w.Wait()
		close(resultch)
	}()

	for listing := range resultch {
		err := repositories.Create(&listing, "scrapped_listings")
		if err != nil {
			continue
		}
	}

	fmt.Println("Finished scrapping aux_predial")
}
