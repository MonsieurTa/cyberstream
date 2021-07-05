package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/MonsieurTa/hypertube/common/db"
	"github.com/MonsieurTa/hypertube/common/entity"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func initEnv() {
	env := os.Getenv("HYPERTUBE_ENV")
	if env == "" {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
}

func main() {
	initEnv()

	db := db.InitDB(&db.PSQLConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Db:       os.Getenv("POSTGRES_DB"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}, &gorm.Config{})

	yesterday := time.Now().AddDate(0, 0, -1)

	basePath := os.Getenv("STATIC_FILES_PATH")
	fmt.Println("Retreiving files in", basePath)

	regulator := NewRegulator(&RegulatorConfig{
		DataDir: []string{basePath, basePath + "/hls"},
		Except:  []string{"hls"},
		Before:  yesterday,
	})

	regulator.Exec()

	gormDB := db.DB().(*gorm.DB)

	err := gormDB.Where("created_at <= ?", yesterday).Delete(&entity.Video{}).Error
	if err != nil {
		panic(err)
	}
}

type Regulator struct {
	dataDir []string
	except  map[string]bool
	before  time.Time
}

type RegulatorConfig struct {
	DataDir []string
	Except  []string
	Before  time.Time
}

func NewRegulator(cfg *RegulatorConfig) *Regulator {
	if cfg == nil {
		panic("no config")
	}
	rv := &Regulator{}
	rv.dataDir = cfg.DataDir
	rv.except = make(map[string]bool)
	rv.before = cfg.Before

	for _, v := range cfg.Except {
		rv.except[v] = true
	}
	return rv
}

func (r *Regulator) Exec() {
	var wg sync.WaitGroup

	for _, dir := range r.dataDir {
		go r.scan(dir, &wg)
		wg.Add(1)
	}
	wg.Wait()
}

func (r *Regulator) scan(dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(dir, "not present")
		return
	}
	for _, f := range files {
		fmt.Println("\tScanning", f.Name())
		if _, present := r.except[f.Name()]; present {
			continue
		}

		if f.ModTime().Before(r.before) {
			filepath := dir + "/" + f.Name()
			fmt.Println("Removing", filepath)
			err = os.RemoveAll(filepath)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Nothing to do")
		}
	}
}
