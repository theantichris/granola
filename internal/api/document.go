package api

const (
	url            = "https://api.granola.ai/v2/get-documents"
	userAgent      = "Granola/5.354.0"
	xClientVersion = "5.354.0"
)

type GranolaResponse struct {
	Documents []Document `json:"docs"`
}

type Document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
