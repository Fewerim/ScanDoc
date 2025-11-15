package typesJSON

type TheInvoice struct {
	DocumentInfo DocumentInfoObj      `json:"document_info"`
	Seller       SellerObj            `json:"seller"`
	Buyer        BuyerObj             `json:"buyer"`
	Shipper      ShipperObj           `json:"shipper"`
	Consignee    ConsigneeInvoiceObj  `json:"consignee"`
	Currency     CurrencyObj          `json:"currency"`
	Items        []ItemInvoiceObj     `json:"items"`
	Totals       TotalsObj            `json:"totals"`
	Signatures   SignaturesInvoiceObj `json:"signatures"`
}

type DocumentInfoObj struct {
	DocumentType string `json:"document_type"`
	Number       string `json:"number"`
	Date         string `json:"date"`
}

type ShipperObj struct {
	SameAsSeller bool    `json:"same_as_seller"`
	Name         *string `json:"name"`
	Address      *string `json:"address"`
	Inn          *string `json:"inn"`
	Kpp          *string `json:"kpp"`
}

type ConsigneeInvoiceObj struct {
	SameAsSeller bool    `json:"same_as_seller"`
	Name         *string `json:"name"`
	Address      *string `json:"address"`
	Inn          *string `json:"inn"`
	Kpp          *string `json:"kpp"`
}

type ItemInvoiceObj struct {
	LineNumber     int     `json:"line_number"`
	ProductName    string  `json:"product_name"`
	UnitCode       string  `json:"unit_code"`
	UnitName       string  `json:"unit_name"`
	Quantity       float64 `json:"quantity"`
	Price          float64 `json:"price"`
	CostWithoutVat float64 `json:"cost_without_vat"`
	VatRate        string  `json:"vat_rate"`
	Vat_amount     float64 `json:"vat_amount"`
	Total_with_vat float64 `json:"total_with_vat"`
}

type TotalsObj struct {
	TotalWithoutVat float64 `json:"total_without_vat"`
	TotalVat        float64 `json:"total_vat"`
	TotalWithVat    float64 `json:"total_with_vat"`
}

type SignaturesInvoiceObj struct {
	Head                   HeadObj                          `json:"head"`
	Accountant             AccountantObj                    `json:"accountant"`
	IndividualEntrepreneur IndividualEntrepreneurInvoiceObj `json:"individual_entrepreneur"`
}

type HeadObj struct {
	Position  string  `json:"position"`
	Name      string  `json:"name"`
	Signature *string `json:"signature"`
}

type AccountantObj struct {
	Position  string  `json:"position"`
	Name      string  `json:"name"`
	Signature *string `json:"signature"`
}

type IndividualEntrepreneurInvoiceObj struct {
	Name      *string `json:"name"`
	Signature *string `json:"signature"`
}
