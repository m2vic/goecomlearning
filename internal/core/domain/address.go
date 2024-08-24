package domain

type Address struct {
	AddressID int    `json:"addressid"`
	House     string `json:"house"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Zip_Code  int    `json:"zip_code"`
}
