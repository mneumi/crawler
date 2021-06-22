package model

type CarDetail struct {
	ID       uint
	BrandId  string
	Name     string
	Price    float64
	ImageURL string
}

type CarModel struct {
	ID         uint
	LogoURL    string
	BrandName  string
	BrandModel string
}
