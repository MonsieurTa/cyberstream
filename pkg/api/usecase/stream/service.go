package stream

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MonsieurTa/hypertube/common/entity"
)

type Service struct {
	streamEndpoint string
	repo           Repository
}

func NewService(repo Repository) UseCase {
	endpoint := `http://` + os.Getenv("MEDIA_HOST") + `:` + os.Getenv("MEDIA_PORT") + `/stream`
	return &Service{endpoint, repo}
}

func (s *Service) Stream(m *entity.Movie) (string, error) {
	storedMovie, err := s.repo.FindByID(m.ID)
	if err != nil {
		return "", err
	}

	if storedMovie != nil {
		return storedMovie.Path, nil
	}

	url, err := stream(s.streamEndpoint, m)
	if err != nil {
		return "", err
	}
	m.Path = url

	_, err = s.repo.Create(m)
	if err != nil {
		return "", err
	}
	return url, nil
}

func stream(endpoint string, m *entity.Movie) (string, error) {
	streamReq := entity.NewStreamRequest(m.Name, m.Magnet)
	data, err := json.Marshal(map[string]interface{}{"stream_request": streamReq})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(b))
	}

	var streamResponse entity.StreamResponse

	err = json.Unmarshal(b, &streamResponse)
	if err != nil {
		return "", err
	}
	return streamResponse.Url, nil
}