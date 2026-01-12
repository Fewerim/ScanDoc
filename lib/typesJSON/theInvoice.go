package typesJSON

import (
	"proWeb/lib/typesJSON/objects"
)

type TheInvoice struct {
	DocumentInfo objects.DocumentInfoObj      `json:"document_info"`
	Seller       objects.SellerObj            `json:"seller"`
	Buyer        objects.BuyerObj             `json:"buyer"`
	Shipper      objects.ShipperObj           `json:"shipper"`
	Consignee    objects.ConsigneeInvoiceObj  `json:"consignee"`
	Currency     objects.CurrencyObj          `json:"currency"`
	Items        []objects.ItemInvoiceObj     `json:"items"`
	Totals       objects.TotalsInvoiceObj     `json:"totals"`
	Signatures   objects.SignaturesInvoiceObj `json:"signatures"`
}
