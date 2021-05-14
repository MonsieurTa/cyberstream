package subsplease

import (
	"net/http"
	"strings"
	"time"
)

type Episode struct {
	Time        string           `json:"time"`
	ReleaseDate dateTime         `json:"release_date"`
	Show        string           `json:"show"`
	Episode     string           `json:"episode"`
	Downloads   []DownloadOption `json:"downloads"`
	Xdcc        string           `json:"xdcc"`
	ImageUrl    string           `json:"image_url"`
	Page        string           `json:"page"`
}

func (e *Episode) HighestResolutionMagnet() string {
	size := len(e.Downloads)
	if size == 0 {
		return ""
	}
	return e.Downloads[size-1].Magnet
}

type DownloadOption struct {
	Res    int    `json:"res,string"`
	Magnet string `json:"magnet"`
}

type subsPlease struct {
	c *http.Client
}

type dateTime time.Time

func (c *dateTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	value = strings.Replace(value, "\\", "", -1) // remove escaped characters

	t, err := time.Parse("01/02/06", value)
	if err != nil {
		return err
	}
	*c = dateTime(t)
	return nil
}

func (c dateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format("01/02/06") + `"`), nil
}

type byReleaseDate []Episode

func (s byReleaseDate) Len() int {
	return len(s)
}

func (s byReleaseDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byReleaseDate) Less(i, j int) bool {
	return time.Time(s[i].ReleaseDate).Before(time.Time(s[j].ReleaseDate))
}

type byResolution []DownloadOption

func (s byResolution) Len() int {
	return len(s)
}

func (s byResolution) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byResolution) Less(i, j int) bool {
	return s[i].Res < s[j].Res
}
