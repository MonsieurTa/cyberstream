package entity

type StreamRequest struct {
	Name   string `json:"name"`
	Magnet string `json:"magnet"`
}

type StreamResponse struct {
	Error string
	Url   string `json:"url"`
}
