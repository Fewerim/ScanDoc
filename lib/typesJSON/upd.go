package typesJSON

import (
	"proWeb/lib/typesJSON/objects"
)

type Upd struct {
	DocumentInfo    objects.DocumentInfoObj    `json:"document_info"`
	Status          int                        `json:"status"`
	Seller          objects.SellerObj          `json:"seller"`
	Consignor       objects.ConsignorObj       `json:"consignor"`
	Consignee       objects.ConsigneeObj       `json:"consignee"`
	PaymentDocument objects.PaymentDocumentObj `json:"paymentDocument"`
	Buyer           objects.BuyerObj           `json:"buyer"`
	Currency        objects.CurrencyObj        `json:"currency"`
	Items           []objects.ItemUpdObj       `json:"items"`
	TotalAmount     objects.TotalAmountObject  `json:"totalAmount"`
	Signatures      objects.SignaturesObj      `json:"signatures"`
	AdditionalInfo  objects.AdditionalInfoObj  `json:"additionalInfo"`
}
