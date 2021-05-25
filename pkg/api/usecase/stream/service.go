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

func (s *Service) Stream(streamReq *entity.StreamRequest) (string, error) {
	storedVideo, err := s.repo.FindByName(streamReq.Name)
	if err != nil {
		return "", err
	}

	if storedVideo != nil {
		return storedVideo.Path, nil
	}

	streamResp, err := stream(s.streamEndpoint, streamReq)
	if err != nil {
		return "", err
	}

	video := entity.NewVideo(
		streamResp.Name,
		streamResp.FileHash,
		streamResp.Url,
		streamReq.Magnet,
	)

	_, err = s.repo.Create(video)
	if err != nil {
		return "", err
	}
	return streamResp.Url, nil
}

func stream(endpoint string, streamReq *entity.StreamRequest) (*entity.StreamResponse, error) {
	data, err := json.Marshal(map[string]interface{}{"stream_request": streamReq})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(b))
	}

	var streamResponse entity.StreamResponse

	err = json.Unmarshal(b, &streamResponse)
	if err != nil {
		return nil, err
	}
	return &streamResponse, nil
}
