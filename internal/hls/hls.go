package hls

import (
	"io"
	"os/exec"
)

func Init(filepath string) (r *io.PipeReader, w *io.PipeWriter, wait func() error) {
	r, w = io.Pipe()
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
	cmd.Stdin = r

	cmd.Start()
	wait = cmd.Wait
	return
}
