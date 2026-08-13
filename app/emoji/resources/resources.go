package resources

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var EmojisMetaData [][]string

func getEmojiID(filename string) (id string) {
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

	if err := os.MkdirAll(filepath.Join(AppDir, "Images"), os.ModePerm); err != nil {
		panic(err)
	}

	if f, err := os.Open(filepath.Join(AppDir, "emojisMetaData.csv")); err != nil && !errors.Is(err, os.ErrNotExist) {
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
			EmojisMetaData = append(EmojisMetaData, line)
		}
		f.Close()
	}

	readAllEmojis()
}

func StoreEmoji(fileName, name, link string, data []byte) {
	id := getEmojiID(fileName)

	emojiMetaData := []string{}
	emojiMetaData = append(emojiMetaData, id)
	emojiMetaData = append(emojiMetaData, name)
	emojiMetaData = append(emojiMetaData, link)
	EmojisMetaData = append(EmojisMetaData, emojiMetaData)

	if f, err := os.Create(filepath.Join(AppDir, "emojisMetaData.csv")); err != nil {
		panic(err)
	} else {
		w := csv.NewWriter(f)
		w.Write([]string{"id", "name", "url"})
		w.WriteAll(EmojisMetaData)
		w.Flush()
		f.Close()
	}

	if f, err := os.Create(filepath.Join(AppDir, "Images", fileName)); err != nil {
		panic(err)
	} else {
		f.Write(data)
		f.Close()
	}
}

type EmojiData struct {
	Id   string
	Name string
	Url  string
	Data []byte
}

var AllEmojisData []EmojiData

func readAllEmojis() {
	path := filepath.Join(AppDir, "Images")
	if files, err := os.ReadDir(path); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				if f, err := os.ReadFile(filepath.Join(path, file.Name())); err != nil {
					fmt.Println(err)
				} else {
					id := getEmojiID(file.Name())
					name := ""
					url := ""
					for _, line := range EmojisMetaData {
						if line[0] == id {
							name = line[1]
							url = line[2]
							break
						}
					}
					if name == "" {
						continue
					}
					AllEmojisData = append(AllEmojisData, EmojiData{id, name, url, f})
				}
			}
		}
	}
}
