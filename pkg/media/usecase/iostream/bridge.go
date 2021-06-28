package iostream

import "fmt"

type Bridge interface {
	GetPlaylist(key string) chan bool
	SetPlaylistReceived(key string) error
	ClosePlaylistRequest(key string)
	PlaylistRequestExists(key string) bool
}

type bridge struct {
	dataDir string

	ongoingRequest map[string]chan bool
}

type BridgeConfig struct {
	DataDir string
}

func NewBridge(cfg *BridgeConfig) Bridge {
	return &bridge{
		dataDir:        cfg.DataDir,
		ongoingRequest: make(map[string]chan bool),
	}
}

func (ios *bridge) GetPlaylist(key string) chan bool {
	ios.ongoingRequest[key] = make(chan bool)
	return ios.ongoingRequest[key]
}

func (ios *bridge) SetPlaylistReceived(key string) error {
	ch, ok := ios.ongoingRequest[key]
	if !ok {
		return fmt.Errorf("bridge: %s request does not exist", key)
	}
	ch <- true
	return nil
}

func (ios *bridge) ClosePlaylistRequest(key string) {
	ch, ok := ios.ongoingRequest[key]
	if !ok {
		return
	}
	close(ch)
	delete(ios.ongoingRequest, key)
}

func (ios *bridge) PlaylistRequestExists(key string) bool {
	_, ok := ios.ongoingRequest[key]
	return ok
}
