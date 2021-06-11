package vtt

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anacrolix/torrent"
)

type VTTConverter interface {
	Convert() []string
}

type vttconverter struct {
	dataDir    string
	tfiles     []*torrent.File
	filesdatas [][]byte
}

func NewVTTConverter(dataDir string, tfiles []*torrent.File) VTTConverter {
	size := len(tfiles)
	filesdatas := make([][]byte, size)
	return &vttconverter{dataDir, tfiles, filesdatas}
}

func (c *vttconverter) Convert() []string {
	var wg sync.WaitGroup

	nbFiles := len(c.tfiles)
	workResults := make(chan string, nbFiles)
	rv := make([]string, 0, nbFiles)
	for _, tfile := range c.tfiles {
		wg.Add(1)
		go c.convert(workResults, tfile, &wg)
	}
	wg.Wait()
	log.Println("DONE CONVERTING SUB")
	for len(workResults) > 0 {
		result := <-workResults
		rv = append(rv, result)
	}
	log.Println("DONE DRAINING SUB")
	close(workResults)
	return rv
}

func (c *vttconverter) convert(results chan string, tfile *torrent.File, wg *sync.WaitGroup) {
	defer wg.Done()

	r := tfile.NewReader()
	_, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err.Error())
		return
	}
	filePath := c.dataDir + "/" + tfile.Path()
	ext := filepath.Ext(tfile.DisplayPath())
	vttName := strings.Replace(tfile.Path(), ext, ".vtt", 1)

	outputFileName := c.dataDir + "/" + vttName
	cmd := exec.Command("ffmpeg", "-i", filePath, outputFileName)

	err = cmd.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err.Error())
		return
	}

	results <- vttName
}
