package scrappers

import (
	"fmt"
	"log"
	"scrapper/internal/repositories"
	"scrapper/internal/structs"
	"scrapper/utils"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func isListingValid(listing structs.ListingItem) bool {
	if strings.Trim(listing.Type, " ") == "" {
		return false
	}
	if strings.Contains(listing.Type, "Galpão") ||
		strings.Contains(listing.Type, "Mônaco") ||
		strings.Contains(listing.Type, "Sobrado") ||
		strings.Contains(listing.Type, "Terreno") {
		return false
	}
	if listing.Image == "" {
		return false
	}
	if !utils.CheckURLImage(listing.Image) {
		return false
	}
	return true
}

func scrapeInitPage() []structs.ListingItem {
	c := colly.NewCollector()

	listings := []structs.ListingItem{}

	c.OnHTML("div.container div.destaquecard__container", func(h *colly.HTMLElement) {
		h.ForEach("a", func(_ int, h *colly.HTMLElement) {
			link := fmt.Sprintf("https://www.casaimoveisimb.com.br%s", h.Attr("href"))
			id, _ := strconv.Atoi(utils.GetIDFromLink(link))
			image := h.ChildAttr("div.destaquecard__img img", "data-src")
			address := h.ChildText("div.destaquecard__img div.row--2 div.destaquecard__image__text p.destaquecard__image__local")
			priceValue := strings.Trim(h.ChildText("div.destaquecard__img div.row--2 div.destaquecard__image__text p.destaquecard__image__valor"), " ")
			price := strings.Replace(priceValue, "R$", "", -1)
			listingType := h.ChildText("div.destaquecard__img div.row--2 div.destaquecard__image__text h2.destaquecard__image__name")
			bedrooms := 0
			bathrooms := 0
			parking := 0
			forSale := h.Attr("data-categoria") == "venda"
			ref := h.ChildText("div.destaquecard__img div.row--1 p")

			h.ForEach("div.destaquecard__img div.row--1 div.destaquecard__img__features", func(i int, h *colly.HTMLElement) {
				h.ForEach("div", func(i int, h *colly.HTMLElement) {
					if i == 0 {
						bedrooms, _ = strconv.Atoi(h.Text)
					}
					if i == 1 {
						bathrooms, _ = strconv.Atoi(h.Text)
					}
					if i == 2 {
						parking, _ = strconv.Atoi(h.Text)
					}
				})
			})

			listing := structs.ListingItem{
				Link:      link,
				Id:        id,
				Image:     image,
				Address:   address,
				Price:     price,
				Type:      listingType,
				Bedrooms:  bedrooms,
				Bathrooms: bathrooms,
				Parking:   parking,
				ForSale:   forSale,
				Ref:       ref,
				Agency:    "casa_imoveis",
			}

			if isListingValid(listing) {
				listings = append(listings, listing)
			}
		})
	})
	c.Visit("https://www.casaimoveisimb.com.br/")
	return listings
}

func scrapePageInside(listing structs.ListingItem, ch chan structs.ListingItem, w *sync.WaitGroup) {
	defer w.Done()

	c := colly.NewCollector()

	c.OnHTML("div.maincontainer div.tags", func(h *colly.HTMLElement) {
		id, _ := strconv.Atoi(utils.GetIDFromLink(h.Response.Request.URL.String()))
		if id == listing.Id {
			h.ForEach("p", func(i int, h *colly.HTMLElement) {
				if strings.Contains(h.Text, "Área Total:") {
					listing.Area = h.ChildText("b")
				}
			})
		}
	})
	c.OnHTML("div.maincontainer div#desc_descricao", func(h *colly.HTMLElement) {
		id, _ := strconv.Atoi(utils.GetIDFromLink(h.Response.Request.URL.String()))
		if id == listing.Id {
			bodyHtml, err := h.DOM.Html()
			if err != nil {
				log.Fatal(err)
			}
			listing.Content = bodyHtml
		}
	})
	c.OnHTML("div.lista-inicial-container div#galeria-inicial div.glide__track ul.glide__slides", func(h *colly.HTMLElement) {
		id, _ := strconv.Atoi(utils.GetIDFromLink(h.Response.Request.URL.String()))
		if id == listing.Id {
			h.ForEach("li", func(i int, h *colly.HTMLElement) {
				imagePhoto := h.ChildAttr("img.item-lista", "src")
				if utils.CheckURLImage(imagePhoto) {
					listing.Photos = append(listing.Photos, structs.Photos{
						Href:          imagePhoto,
						ListingItemId: listing.Id,
					})
				}
			})
		}
	})
	c.Visit(listing.Link)
	ch <- listing
}

func CasaImoveisExecute(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Starting scrapping casa_imoveis...")

	listings := scrapeInitPage()

	listingsChannel := make(chan structs.ListingItem)

	var w sync.WaitGroup

	for _, listing := range listings {
		w.Add(2)
		go scrapePageInside(listing, listingsChannel, &w)
		go scrapePageInside(listing, listingsChannel, &w)
	}

	go func() {
		w.Wait()
		close(listingsChannel)
	}()

	for listing := range listingsChannel {
		err := repositories.Create(&listing, "scrapped_listings")
		if err != nil {
			continue
		}
	}

	fmt.Println("Finished scrapping casa_imoveis")
}
