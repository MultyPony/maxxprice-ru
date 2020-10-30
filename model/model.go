package model

type Flat struct {
	ID              string   `bson:"_id"`
	Price           int64    `bson:"Price"`
	LocationName    string   `bson:"LocationName"`
	Square          int64    `bson:"Square"`
	PerSquare       int64    `bson:"PerSquare"`
	SortTimeStamp   int64    `bson:"SortTimeStamp"`
	Title           string   `bson:"Title"`
	Agency          string   `bson:"Agency"`
	RawPhone        string   `bson:"RawPhone"`
	Link            string   `bson:"Link"`
	SplittedAddress []string `bson:"SplittedAddress"`
}

type CallReq struct {
	ID           int64  `bson:"_id"`
	CreationTime int64  `bson:"CreationTime"`
	Square       int64  `bson:"Square"`
	Data         string `bson:"Data"`
	Street       string `bson:"Street"`
	Response1    int64  `bson:"Response1"`
	Response2    int64  `bson:"Response2"`
}

type AddressDetailed struct {
	Metro    string `bson:"Metro"`
	District string `bson:"District"`
	Road     string `bson:"Road"`
	Name     string `bson:"Name"`
}

type Location struct {
	ID   int64  `bson:"_id"`
	Name string `bson:"Name"`
}

type AvitoFlatResponse struct {
	LastStamp int64         `json:"lastStamp"`
	Title     string        `json:"title"`
	Context   string        `json:"context"`
	Items     []ArticleItem `json:"items"`
}

type ArticleItem struct {
	ID                    int64                   `json:"id"`
	CaregoryID            int64                   `json:"categoryId"`
	LocationID            int64                   `json:"locationId"`
	Title                 string                  `json:"title"`
	ImagesCount           int                     `json:"imagesCount"`
	IsActive              bool                    `json:"isActive"`
	URLPath               string                  `json:"urlPath"`
	SortTimestamp         int64                   `json:"sortTimeStamp"`
	Category              AvitoCategory           `json:"category"`
	Location              AvitoLocation           `json:"location"`
	AddressDetailed       AvitoAddressDetailed    `json:"addressDetailed"`
	HasViseo              bool                    `json:"hasVideo"`
	AvitoListingImageURLs AvitoListingImageURLs   `json:"listingImageUrls"`
	Images                []AvitoListingImageURLs `json:"images"`
	PriceString           string                  `json:"priceString"`
	IsMarketplace         bool                    `json:"isMarketplace"`
	IsFavorite            bool                    `json:"IsFavorite"`
	ImagesAlt             string                  `json:"imagesAlt"`
	PriceDetailed         AvitoPriceDetailed      `json:"priceDetailed"`
	GeoForItems           GeoItem                 `json:"geoForItems"`
}

type AvitoPriceDetailed struct {
	Title                     AvitoTitle `json:"title"`
	TitleDative               string     `json:"titleDative"`
	Postfix                   string     `json:"postfix"`
	Enabled                   bool       `json:"enabled"`
	PostfixShort              string     `json:"postfixShort"`
	WasLowered                bool       `json:"wasLowered"`
	HasValue                  bool       `json:"hasValue"`
	String                    string     `json:"string"`
	FullString                string     `json:"fullString"`
	WalueWithoutDiscoutSigned string     `json:"valueWithoutDiscountSigned"`
	Value                     int64      `json:"value"`
	MinPrice                  int64      `json:"minPrice"`
	MinPriceUpdate            string     `json:"minPriceUpdate"`
	Metric                    string     `json:"metric"`
	Hint                      string     `json:"hint"`
}

type GeoItem struct {
	ItemID           int64           `json:"itemId"`
	FormattedAddress string          `json:"formattedAddress"`
	GeoReferences    []GeoReferences `json:"geoReferences"`
}

type GeoReferences struct {
	Content string `json:"content"`
}

type AvitoTitle struct {
	Full  string `json:"full"`
	Short string `json:"short"`
}

type AvitoListingImageURLs struct {
	Miniature  string `json:"208x156"`
	Catalog    string `json:"catalog"`
	CatalogVIP string `json:"catalog_vip"`
}

type AvitoAddressDetailed struct {
	MetroName    string `json:"metroName"`
	DistrictName string `json:"districtName"`
	RoadName     string `json:"roadName"`
	LocationName string `json:"locationName"`
}

type AvitoLocation struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	NamePrepositional string `json:"namePrepositional"`
	NameGenitive      string `json:"nameGenitive"`
}

type AvitoCategory struct {
	ID       int64  `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	ParentID int64  `json:"parentId"`
}
