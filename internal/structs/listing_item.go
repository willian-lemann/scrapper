package structs

type ListingItem struct {
	Id               int      `json:"id,omitempty"`
	Link             string   `json:"link"`
	Image            string   `json:"image"`
	Address          string   `json:"address"`
	Price            string   `json:"price"`
	Area             string   `json:"area"`
	Bedrooms         int      `json:"bedrooms"`
	Bathrooms        int      `json:"bathrooms"`
	Type             string   `json:"type"`
	ForSale          bool     `json:"forSale"`
	Parking          int      `json:"parking"`
	Content          string   `json:"content"`
	Photos           []Photos `json:"photos"`
	Agency           string   `json:"agency"`
	PlaceholderImage string   `json:"placeholderImage"`
	Ref              string   `json:"ref"`
}

type Photos struct {
	ListingItemId int    `json:"listingItemId"`
	Href          string `json:"href"`
}

func NewListingItem(listing ListingItem) *ListingItem {
	return &ListingItem{
		Id:               listing.Id,
		Link:             listing.Link,
		Image:            listing.Image,
		Address:          listing.Address,
		Price:            listing.Price,
		Area:             listing.Area,
		Bedrooms:         listing.Bedrooms,
		Bathrooms:        listing.Bathrooms,
		Type:             listing.Type,
		ForSale:          listing.ForSale,
		Parking:          listing.Parking,
		Content:          listing.Content,
		Photos:           listing.Photos,
		Agency:           listing.Agency,
		PlaceholderImage: listing.PlaceholderImage,
		Ref:              listing.Ref,
	}
}

func (l *ListingItem) CreateListingWithEmptyId(listing ListingItem) ListingItem {
	return ListingItem{
		Link:             listing.Link,
		Ref:              listing.Ref,
		Image:            listing.Image,
		Address:          listing.Address,
		Price:            listing.Price,
		Area:             listing.Area,
		Bedrooms:         listing.Bedrooms,
		Bathrooms:        listing.Bathrooms,
		Type:             listing.Type,
		ForSale:          listing.ForSale,
		Parking:          listing.Parking,
		Content:          listing.Content,
		Photos:           listing.Photos,
		Agency:           listing.Agency,
		PlaceholderImage: listing.PlaceholderImage,
	}
}
