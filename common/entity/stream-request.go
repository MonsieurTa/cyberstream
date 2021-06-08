package entity

type StreamRequest struct {
	InfoHash string `form:"info_hash" json:"info_hash" binding:"required"`
	Magnet   string `form:"magnet" json:"magnet" binding:"required"`
}

type StreamResponse struct {
	Name          string   `json:"name,omitempty"`
	Ext           string   `json:"ext,omitempty"`
	InfoHash      string   `json:"info_hash,omitempty"`
	MediaURL      string   `json:"media_url,omitempty"`
	SubtitlesURLs []string `json:"subtitles_urls,omitempty"`
	Error         string   `json:"error,omitempty"`
}
