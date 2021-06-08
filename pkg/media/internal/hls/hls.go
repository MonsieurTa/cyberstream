package hls

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/MonsieurTa/hypertube/common/file"
	"github.com/anacrolix/torrent"
)

type HLSConverter interface {
	Convert()
	WaitUntilReady() error
	WaitUntilDone()
	Close()
	PlaylistPath() string
}

type hlsConverter struct {
	cfg *Config

	input  *io.PipeWriter
	ffmpeg *exec.Cmd

	ready   chan bool
	done    chan bool
	errChan chan error
}

type Config struct {
	StreamFile *torrent.File
	OutputDir  string
}

func NewHLSConverter(cfg *Config) HLSConverter {
	cmd := initFfmpeg(cfg.OutputDir)

	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdin = pipeReader
	return &hlsConverter{
		cfg:     cfg,
		input:   pipeWriter,
		ffmpeg:  cmd,
		ready:   make(chan bool),
		done:    make(chan bool),
		errChan: make(chan error),
	}
}

func (c *hlsConverter) PlaylistPath() string {
	output := c.cfg.OutputDir + "/master.m3u8"
	return output
}

func (c *hlsConverter) WaitUntilReady() error {
	for {
		select {
		case err := <-c.errChan:
			return err
		case <-c.ready:
			return nil
		default:
		}
		time.Sleep(time.Microsecond * 100)
	}
}

func (c *hlsConverter) WaitUntilDone() {
	<-c.done
}

func (c *hlsConverter) Close() {
	close(c.ready)
	close(c.done)
	close(c.errChan)
}

func initFfmpeg(outputDir string) *exec.Cmd {
	input := []string{
		"-i", "pipe:0",
	}
	codecs := []string{
		"-preset", "veryfast",
		"-pix_fmt", "yuv420p",
		"-crf", "21",
		"-c:v", "libx264", // video codec
		"-c:a", "aac", // audio codec
		"-b:a", "128k", // audio bitrate
		"-ac", "2", // number of audio channels
		"-c:s", "webvtt", // subtitle codec
	}
	streamMap := []string{
		"-map", "0:v",
		"-map", "0:a",
		"-map", "0:s",
		"-var_stream_map", "v:0,a:0,s:0,sgroup:subtitle",
	}
	hls := []string{
		"-f", "hls",
		"-hls_time", "15",
		"-hls_list_size", "0",
		"-hls_playlist_type", "event",
		"-master_pl_name", "master.m3u8",
		outputDir + "/out_%v.m3u8",
	}

	input = append(input, codecs...)
	input = append(input, streamMap...)
	input = append(input, hls...)

	cmd := exec.Command("ffmpeg", input...)
	return cmd
}

func (c *hlsConverter) Convert() {
	stopGuard := make(chan bool)

	os.MkdirAll(c.cfg.OutputDir, os.ModePerm)

	err := c.ffmpeg.Start()
	if err != nil {
		c.errChan <- err
		stopGuard <- true
		return
	}

	go c.convert()

	// ready when hls out.m3u8 is created
	playlistPath := c.cfg.OutputDir + "/master.m3u8"
	go guard(playlistPath, stopGuard, c.ready)
}

func guard(playlistPath string, stop, ready chan bool) {
	for {
		select {
		case <-stop:
			return
		default:
		}
		if file.Exists(playlistPath) {
			ready <- true
			return
		}
		time.Sleep(time.Microsecond * 100)
	}
}

func (c *hlsConverter) convert() {
	go listenError(c.ffmpeg, c.errChan)
	defer func() { c.done <- true }()

	r := c.cfg.StreamFile.NewReader()
	defer r.Close()

	buf := make([]byte, c.cfg.StreamFile.Torrent().Info().PieceLength)
	at := int64(0)
	end := c.cfg.StreamFile.Length()

	r.SetReadahead(end / 100 * 5)
	for at < end {
		// reading from the torrent.Reader will download the resource asked
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			c.errChan <- err
			return
		}

		_, err = c.input.Write(buf)
		if err != nil && err != io.ErrClosedPipe {
			c.errChan <- err
			return
		}
		at += int64(n)
	}
}

func listenError(cmd *exec.Cmd, c chan error) {
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			c <- err
		}
	}
}
