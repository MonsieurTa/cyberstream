package entity

type StreamRequest struct {
	Name   string `form:"name" json:"name"`
	Magnet string `form:"magnet" json:"magnet"`
}

type StreamResponse struct {
	Name     string `json:"name"`
	FileHash string `json:"file_hash"`
	Url      string `json:"url"`
}

func NewStreamRequest(name, magnet string) *StreamRequest {
	return &StreamRequest{name, magnet}
}

func NewStreamResponse(name, fileHash, url string) *StreamResponse {
	return &StreamResponse{name, fileHash, url}
}
