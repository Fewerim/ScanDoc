package handlers

type ProcessRequest struct {
	Format string `json:"format"`
	Data   []byte `json:"data"`
}

type ProcessResult struct {
	DocumentType string  `json:"document_type"`
	JsonResult   string  `json:"json_result"`
	Progress     float64 `json:"progress"`
	Error        string  `json:"error,omitempty"`
}
