package resources

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var GifsMetaData [][]string

func getGifID(filename string) (id string) {
	splited := strings.Split(filename, ".")
	id = strings.Join(splited[:len(splited)-1], ".")
	return
}

var AppDir string

func Init() {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	AppDir = filepath.Join(dir, "EmojiPicker")

	if err := os.MkdirAll(filepath.Join(AppDir, "Gifs"), os.ModePerm); err != nil {
		panic(err)
	}

	if f, err := os.Open(filepath.Join(AppDir, "gifsMetaData.csv")); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	} else if err == nil {
		r := csv.NewReader(f)
		data, err := r.ReadAll()
		if err != nil {
			panic(err)
		}
		for i, line := range data {
			if i == 0 {
				continue
			}
			GifsMetaData = append(GifsMetaData, line)
		}
		f.Close()
	}

	readAllGifs()
}

func StoreGif(fileName, link string, data []byte) {
	id := getGifID(fileName)

	gifMetaData := []string{}
	gifMetaData = append(gifMetaData, id)
	gifMetaData = append(gifMetaData, link)
	GifsMetaData = append(GifsMetaData, gifMetaData)

	if f, err := os.Create(filepath.Join(AppDir, "gifsMetaData.csv")); err != nil {
		panic(err)
	} else {
		w := csv.NewWriter(f)
		w.Write([]string{"id", "url"})
		w.WriteAll(GifsMetaData)
		w.Flush()
		f.Close()
	}

	if f, err := os.Create(filepath.Join(AppDir, "Gifs", fileName)); err != nil {
		panic(err)
	} else {
		f.Write(data)
		f.Close()
	}
}

type GifData struct {
	Id   string
	Url  string
	Data []byte
}

var AllGifsData []GifData

func readAllGifs() {
	path := filepath.Join(AppDir, "Gifs")
	if files, err := os.ReadDir(path); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				if f, err := os.ReadFile(filepath.Join(path, file.Name())); err != nil {
					fmt.Println(err)
				} else {
					id := getGifID(file.Name())
					found := false
					url := ""
					for _, line := range GifsMetaData {
						if line[0] == id {
							url = line[1]
							found = true
							break
						}
					}
					if !found {
						continue
					}
					AllGifsData = append(AllGifsData, GifData{id, url, f})
				}
			}
		}
	}
}
