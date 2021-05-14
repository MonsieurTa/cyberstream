package entity

type StreamRequest struct {
	Name   string `json:"name"`
	Magnet string `json:"magnet"`
}

type StreamResponse struct {
	Error string `json:"error"`
	Url   string `json:"url"`
}

func NewStreamRequest(name, magnet string) *StreamRequest {
	return &StreamRequest{name, magnet}
}

func NewStreamResponse(url string) *StreamResponse {
	return &StreamResponse{Url: url}
}
