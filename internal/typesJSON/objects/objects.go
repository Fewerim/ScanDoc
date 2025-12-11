package objects

type DocumentInfoObj struct {
	DocumentType string `json:"document_type"`
	Number       string `json:"number"`
	Date         string `json:"date"`
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

type ItemUpdObj struct {
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
	VatAmount      float64 `json:"vat_amount"`
	TotalWithVat   float64 `json:"total_with_vat"`
}

type TotalsInvoiceObj struct {
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

type ContractObj struct {
	Number string `json:"number"`
	Date   string `json:"date"`
}

type CompanyObj struct {
	Name                 string `json:"name"`
	Inn                  string `json:"inn"`
	Kpp                  string `json:"kpp"`
	Address              string `json:"address"`
	BankAccount          string `json:"bank_account"`
	BankName             string `json:"bank_name"`
	Bik                  string `json:"bik"`
	CorrespondentAccount string `json:"correspondent_account"`
}

type ItemTorgObj struct {
	LineNumber       int     `json:"line_number"`
	ProductName      string  `json:"product_name"`
	ProductCode      string  `json:"product_code"`
	Unit             string  `json:"unit"`
	OkeiCode         string  `json:"okei_code"`
	PackageType      string  `json:"package_type"`
	UnitsPerPlace    int     `json:"units_per_place"`
	PlacesCount      int     `json:"places_count"`
	GrossWeight      float64 `json:"gross_weight"`
	NetWeight        float64 `json:"net_weight"`
	Price            float64 `json:"price"`
	AmountWithoutVat float64 `json:"amount_without_vat"`
	VatRatePercent   int     `json:"vat_rate_percent"`
	VatAmount        float64 `json:"vat_amount"`
	AmountWithVat    float64 `json:"amount_with_vat"`
}

type TotalsTorgObj struct {
	TotalQuantity         int     `json:"total_quantity"`
	TotalGrossWeight      float64 `json:"total_gross_weight"`
	TotalNetWeight        float64 `json:"total_net_weight"`
	TotalAmountWithoutVat float64 `json:"total_amount_without_vat"`
	TotalVatAmount        float64 `json:"total_vat_amount"`
	TotalAmountWithVat    float64 `json:"total_amount_with_vat"`
}

type AttachmentsTorgObj struct {
	AttachmentsDescription string `json:"attachment_description"`
	AttachmentsPagesCount  int    `json:"attachments_pages_count"`
}

type PowerOfAttorneyTorgObj struct {
	Number string `json:"number"`
	Date   string `json:"date"`
}

type SignaturesTorgObj struct {
	ReleasedByPosition       string `json:"releasedByPosition"`
	ReleasedBySignature      string `json:"releasedBySignature"`
	ReleasedByName           string `json:"releasedByName"`
	ChiefAccountantSignature string `json:"chiefAccountantSignature"`
	ChiefAccountantName      string `json:"chiefAccountantName"`
	CargoReceivedByPosition  string `json:"cargoReceivedByPosition"`
	CargoReceivedBySignature string `json:"cargoReceivedBySignature"`
	CargoReceivedByName      string `json:"cargoReceivedByName"`
	CargoReleasedByPosition  string `json:"cargoReleasedByPosition"`
	CargoReleasedBySignature string `json:"cargoReleasedBySignature"`
	CargoReleasedByName      string `json:"cargoReleasedByName"`
	RecipientPosition        string `json:"recipientPosition"`
	RecipientSignature       string `json:"recipientSignature"`
	RecipientName            string `json:"recipientName"`
}

type StampTorgObj struct {
	Description string `json:"description"`
	Date        string `json:"date"`
}
