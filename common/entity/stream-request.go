package entity

type StreamRequest struct {
	Magnet string `json:"magnet"`
}

type StreamResponse struct {
	Error string
	Url   string `json:"url"`
}
