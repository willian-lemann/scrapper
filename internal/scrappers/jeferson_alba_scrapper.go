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

func getLinks() []string {
	var links []string

	for i := 1; i <= 3; i++ {
		link := "https://imobiliariajefersonealba.com.br/alugueis/pesquisa/todos/imbituba/todos/todos/"
		links = append(links, fmt.Sprintf("%s%d", link, i))
	}

	for i := 1; i <= 5; i++ {
		link := "https://imobiliariajefersonealba.com.br/vendas/pesquisa/apartamento/imbituba/todos/todos/"
		links = append(links, fmt.Sprintf("%s%d", link, i))
	}

	return links
}

func getListingItems(linkToVisit string) []structs.ListingItem {
	c := colly.NewCollector()
	listings := []structs.ListingItem{}

	c.OnHTML("div.row div.col-imovel", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a", "href")
		id, _ := strconv.Atoi(utils.GetIDFromLink(link))
		image := h.ChildAttr("div.box-content div.box-imovel-image img", "src")
		forSale := h.ChildText("div.box-content div.box-imovel-infos div.box-imovel-tag span") == "Venda"
		price := strings.Split(h.ChildText("div.box-content div.box-imovel-infos span.--price"), " ")[1]
		listingType := h.ChildText("div.box-content div.box-imovel-infos span.--type")
		address := h.ChildText("div.box-content div.box-imovel-infos span.--location")
		area := ""
		bedrooms := 0
		bathrooms := 0
		parking := 0
		ref := utils.GetIDFromLink(link)
		h.ForEach("div.box-content div.box-imovel-infos ul.box-imovel-items li", func(i int, h *colly.HTMLElement) {
			if strings.Contains(h.Text, "vagas") {
				parkingQuantity, _ := strconv.Atoi(h.ChildText("strong"))
				parking = parkingQuantity
			}

			if strings.Contains(h.Text, "dormitório") {
				bedroomsQuantity, _ := strconv.Atoi(h.ChildText("strong"))
				bedrooms = bedroomsQuantity
			}

			if strings.Contains(h.Text, "banheiro") {
				bathroomsQuantity, _ := strconv.Atoi(h.ChildText("strong"))
				bathrooms = bathroomsQuantity
			}

			if strings.Contains(h.Text, "m²") {
				area = h.ChildText("strong")
			}
		})

		listing := structs.ListingItem{
			Id:        id,
			Link:      link,
			Image:     image,
			Address:   address,
			Price:     price,
			Area:      area,
			Bedrooms:  bedrooms,
			Bathrooms: bathrooms,
			Type:      listingType,
			ForSale:   forSale,
			Parking:   parking,
			Ref:       ref,
			Agency:    "jeferson_alba",
		}

		listings = append(listings, listing)
	})

	c.Visit(linkToVisit)
	return listings
}

func getListingItem(listing structs.ListingItem) structs.ListingItem {
	c := colly.NewCollector()
	c.OnHTML("div.row div div.imovel-content-section", func(h *colly.HTMLElement) {
		id, _ := strconv.Atoi(utils.GetIDFromLink(h.Response.Request.URL.String()))
		if id == listing.Id {
			bodyHtml, err := h.DOM.Html()
			if err != nil {
				log.Fatal(err)
			}
			listing.Content = bodyHtml
		}
	})
	c.OnHTML("div#imovel-fotos div.container div.img-gallery-magnific", func(h *colly.HTMLElement) {
		h.ForEach("div.magnific-img", func(i int, h *colly.HTMLElement) {
			url := h.Response.Request.URL.String()
			id, _ := strconv.Atoi(utils.GetIDFromLink(url))
			srcImage := h.ChildAttr("a img", "src")

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

	})

	c.Visit(listing.Link)
	return listing
}

func worker(link string, ch chan structs.ListingItem, w *sync.WaitGroup) {
	defer w.Done()

	listings := getListingItems(link)

	for _, listing := range listings {
		newListing := getListingItem(listing)
		if len(newListing.Photos) == 0 {
			continue
		}
		ch <- newListing
	}
}

func ExecuteJefersonAlba(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Starting scrapping jeferson_alba...")

	links := getLinks()

	resultch := make(chan structs.ListingItem)

	var w sync.WaitGroup

	for _, link := range links {
		w.Add(1)
		go worker(link, resultch, &w)
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

	fmt.Println("Finished scrapping jeferson_alba")
}
