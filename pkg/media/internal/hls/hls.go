package hls

import (
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/anacrolix/torrent"
)

type HLSConverter interface {
	Convert() error
	WaitUntilReady()
	WaitUntilDone()
	Close() error
}

type hlsConverter struct {
	input   *io.PipeWriter
	ffmpeg  *exec.Cmd
	torrent *torrent.Torrent

	dataDst string

	progress chan int64
	ready    chan bool
	done     chan bool
}

func NewHLSConverter(dst string, t *torrent.Torrent) HLSConverter {
	pipeReader, pipeWriter := io.Pipe()

	cmd := initFfmpeg(dst)
	cmd.Stdin = pipeReader
	return &hlsConverter{
		input:    pipeWriter,
		dataDst:  dst,
		ffmpeg:   cmd,
		torrent:  t,
		progress: make(chan int64),
		ready:    make(chan bool),
		done:     make(chan bool),
	}
}

func (c *hlsConverter) WaitUntilReady() {
	<-c.ready
}

func (c *hlsConverter) WaitUntilDone() {
	<-c.done
}

func (c *hlsConverter) Close() error {
	close(c.progress)
	close(c.ready)
	close(c.done)
	return c.ffmpeg.Wait()
}

func initFfmpeg(filepath string) *exec.Cmd {
	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0",
		"-c:v", "libx264", "-crf", "21", "-preset", "veryfast",
		"-c:a", "aac", "-b:a", "128k", "-ac", "2",
		"-f", "hls",
		"-hls_time", "15",
		"-hls_playlist_type", "event",
		filepath,
	)
	return cmd
}

func (c *hlsConverter) Convert() error {
	err := c.ffmpeg.Start()
	if err != nil {
		return err
	}
	go c.convert()

	// ready when 1% downloaded
	go func() {
		threshold := c.torrent.Length() / 100
		for at := range c.progress {
			if at >= threshold {
				break
			}
			time.Sleep(time.Microsecond * 100)
		}
		c.ready <- true
	}()
	return nil
}

func (c *hlsConverter) convert() {
	r := c.torrent.NewReader()
	buf := make([]byte, c.torrent.Info().PieceLength)
	at := int64(0)
	end := c.torrent.Length()

	r.SetReadahead(end / 100 * 5)
	for at < end {
		// reading from the torrent.Reader will download the resource asked
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		_, err = c.input.Write(buf)
		if err != nil {
			log.Fatal(err)
		}
		at += int64(n)
		select {
		case c.progress <- at:
		default:
		}
	}
	c.done <- true

	r.Close()
}
