package hls

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

const TCP_MAX_BUFFER_SIZE = 8096000

type HLSConverter interface {
	Convert()
	WaitUntilDone()
	Close()
	PlaylistPath() string
}

type hlsConverter struct {
	cfg *Config

	input  *io.PipeWriter
	ffmpeg *exec.Cmd

	done    chan bool
	errChan chan error
}

type Config struct {
	Reader      io.Reader
	Length      int64
	PieceLength int64
	OutputDir   string
}

func NewHLSConverter(cfg *Config) HLSConverter {
	cmd := initFfmpeg(cfg.OutputDir)

	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdin = pipeReader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return &hlsConverter{
		cfg:     cfg,
		input:   pipeWriter,
		ffmpeg:  cmd,
		done:    make(chan bool),
		errChan: make(chan error),
	}
}

func (c *hlsConverter) PlaylistPath() string {
	output := c.cfg.OutputDir + "/master.m3u8"
	return output
}

func (c *hlsConverter) WaitUntilDone() {
	<-c.done
}

func (c *hlsConverter) Close() {
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
	url := fmt.Sprintf("%s/hls/%s/out_%%v.m3u8", os.Getenv("MEDIA_PRIVATE_URL"), outputDir)
	hls := []string{
		"-f", "hls",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_playlist_type", "event",
		"-master_pl_name", "master.m3u8",
		"-method", "POST",
		url,
	}

	input = append(input, codecs...)
	input = append(input, streamMap...)
	input = append(input, hls...)

	return exec.Command("ffmpeg", input...)
}

func (c *hlsConverter) Convert() {
	err := c.ffmpeg.Start()
	if err != nil {
		c.errChan <- err
		return
	}

	go c.convert()
}

func (c *hlsConverter) convert() {
	go listenError(c.ffmpeg, c.errChan)
	defer func() {
		c.done <- true
		c.input.Close()
	}()

	r := c.cfg.Reader

	buf := make([]byte, TCP_MAX_BUFFER_SIZE)
	at := int64(0)
	end := c.cfg.Length

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
		}
		c <- err
		cmd.Process.Kill()
	}
}
