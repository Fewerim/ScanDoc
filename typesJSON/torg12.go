package typesJSON

import "proWeb/typesJSON/objects"

type Torg12 struct {
	DocumentInfo objects.DocumentInfoObj        `json:"document_info"`
	Contract     objects.ContractObj            `json:"contract"`
	Sender       objects.CompanyObj             `json:"sender"`
	Receiver     objects.CompanyObj             `json:"receiver"`
	Supplier     objects.CompanyObj             `json:"supplier"`
	Payer        objects.CompanyObj             `json:"payer"`
	Items        []objects.ItemTorgObj          `json:"items"`
	Totals       objects.TotalsTorgObj          `json:"totals"`
	Attachments  objects.AttachmentsTorgObj     `json:"attachments"`
	Attorney     objects.PowerOfAttorneyTorgObj `json:"attorney"`
	Signatures   objects.SignaturesTorgObj      `json:"signatures"`
	Stamps       []objects.StampTorgObj         `json:"stamps"`
}
