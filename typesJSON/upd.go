package typesJSON

type Upd struct {
	DocumentNumber  string             `json:"documentNumber"`
	DocumentDate    string             `json:"documentDate"`
	Status          int                `json:"status"`
	Seller          SellerObj          `json:"seller"`
	Consignor       ConsignorObj       `json:"consignor"`
	Consignee       ConsigneeObj       `json:"consignee"`
	PaymentDocument PaymentDocumentObj `json:"paymentDocument"`
	Buyer           BuyerObj           `json:"buyer"`
	Currency        CurrencyObj        `json:"currency"`
	Items           []ItemObj          `json:"items"`
	TotalAmount     TotalAmountObject  `json:"totalAmount"`
	Signatures      SignaturesObj      `json:"signatures"`
	AdditionalInfo  AdditionalInfoObj  `json:"additionalInfo"`
}

type SellerObj struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Inn     string `json:"inn"`
	Kpp     string `json:"kpp"`
}

type ConsignorObj struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ConsigneeObj struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type PaymentDocumentObj struct {
	Number string `json:"number"`
	Date   string `json:"date"`
}

type BuyerObj struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Inn     string `json:"inn"`
	Kpp     string `json:"kpp"`
}

type CurrencyObj struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type ItemObj struct {
	LineNumber              int     `json:"lineNumber"`
	ProductCode             string  `json:"productCode"`
	ProductName             string  `json:"productName"`
	Unit                    string  `json:"unit"`
	Quantity                float64 `json:"quantity"`
	Price                   float64 `json:"price"`
	AmountWithoutVAT        float64 `json:"amountWithoutVAT"`
	ExciseAmount            float64 `json:"exciseAmount"`
	VatRatePercent          int     `json:"vatRatePercent"`
	VatAmount               float64 `json:"vatAmount"`
	AmountWithVAT           float64 `json:"amountWithVAT"`
	CountryOfOrigin         string  `json:"countryOfOrigin"`
	CustomDeclarationNumber string  `json:"customDeclarationNumber"`
	CustomCode              string  `json:"customCode"`
	ShortName               string  `json:"shortName"`
}

type TotalAmountObject struct {
	TotalWithoutVAT float64 `json:"totalWithoutVAT"`
	TotalVAT        float64 `json:"totalVAT"`
	TotalWithVAT    float64 `json:"totalWithVAT"`
}

type SignaturesObj struct {
	AuthorizedPerson       AuthorizedPersonObj       `json:"authorizedPerson"`
	ChiefAccountant        ChiefAccountantObj        `json:"chiefAccountant"`
	IndividualEntrepreneur IndividualEntrepreneurObj `json:"individualEntrepreneur"`
}

type AuthorizedPersonObj struct {
	Position string `json:"position"`
	Name     string `json:"name"`
}
type ChiefAccountantObj struct {
	Name string `json:"name"`
}
type IndividualEntrepreneurObj struct {
	Name                string `json:"name"`
	RegistrationDetails string `json:"registrationDetails"`
}

type AdditionalInfoObj struct {
	ShipmentBasis      string                 `json:"shipmentBasis"`
	TransportDocuments string                 `json:"transportDocuments"`
	Notes              string                 `json:"notes"`
	ResponsiblePerson  []ResponsiblePersonObj `json:"responsiblePerson"`
	DocumentCreator    DocumentCreatorObj     `json:"documentCreator"`
}
type ResponsiblePersonObj struct {
	Position string `json:"position"`
	Name     string `json:"name"`
}
type DocumentCreatorObj struct {
	Name string `json:"name"`
	Inn  string `json:"inn"`
	Kpp  string `json:"kpp"`
}
