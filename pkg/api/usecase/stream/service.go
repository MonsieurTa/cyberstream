package stream

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MonsieurTa/hypertube/common/entity"
)

type Service struct {
	streamEndpoint string
	repo           Repository
}

func NewService(repo Repository) UseCase {
	endpoint := os.Getenv("MEDIA_INTERNAL_URL") + ":" + os.Getenv("MEDIA_PORT") + `/stream`
	return &Service{endpoint, repo}
}

func (s *Service) Stream(streamReq *entity.StreamRequest) (*entity.StreamResponse, error) {
	storedVideo, err := s.repo.FindByHash(streamReq.InfoHash)
	if err != nil {
		return nil, err
	}

	if storedVideo != nil {
		if err != nil {
			return nil, err
		}
		return &entity.StreamResponse{
			Name:          storedVideo.Name,
			Ext:           filepath.Ext(storedVideo.Name),
			InfoHash:      storedVideo.Hash,
			MediaURL:      storedVideo.FileURL,
			SubtitlesURLs: storedVideo.SubtitlesURLs,
		}, nil
	}

	streamResp, err := stream(s.streamEndpoint, streamReq)
	if err != nil {
		return nil, err
	}

	video := entity.NewVideo(
		streamResp.Name,
		streamResp.InfoHash,
		streamResp.MediaURL,
		streamReq.Magnet,
		streamResp.SubtitlesURLs,
	)

	_, err = s.repo.Create(video)
	if err != nil {
		return nil, err
	}
	return streamResp, nil
}

func stream(endpoint string, streamReq *entity.StreamRequest) (*entity.StreamResponse, error) {
	data, err := json.Marshal(streamReq)
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

	var streamResponse entity.StreamResponse

	err = json.Unmarshal(b, &streamResponse)
	if err != nil {
		return nil, err
	}
	if len(streamResponse.Error) > 0 {
		return nil, errors.New(streamResponse.Error)
	}
	return &streamResponse, nil
}
