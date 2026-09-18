package emoji

import (
	"EmojiPicker/app/emoji/resources"
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/webp"

	"github.com/AllenDang/cimgui-go/backend"
)

func downloadFile(url string) ([]byte, error) {
	url = strings.Split(url, "?")[0]
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status code: %s", resp.Status)
	}

	out := bytes.NewBuffer(nil)

	io.Copy(out, resp.Body)
	return out.Bytes(), nil
}

type Emoji struct {
	Name     string
	Url      string
	Keywords []string
	Img      *image.RGBA
}

var AllEmojis []Emoji

func Init() {
	resources.Init()
	es := resources.AllEmojisData
	for _, e := range es {
		imgImage, _, err := image.Decode(bytes.NewReader(e.Data))
		if err != nil {
			fmt.Println(err)
			continue
		}
		AllEmojis = append(AllEmojis, Emoji{e.Name, e.Url, []string{}, backend.ImageToRgba(imgImage)})
	}
}

func GetAllEMojis() []Emoji {
	return AllEmojis
}

func getFileName(url string) (name string) {
	name = filepath.Base(url)
	splited := strings.Split(name, "?")
	name = splited[0]
	return
}

func SaveEmoji(url, name string) (e Emoji, finalerr error) {
	url = strings.Split(url, "?")[0]
	dirs := strings.Split(url, "/")
	fmt.Println(url)
	fmt.Println(dirs)
	if len(dirs) < 2 || dirs[2] != "cdn.discordapp.com" {
		finalerr = fmt.Errorf("%s não é uma URL para o discord", url)
		return
	}
	data, err := downloadFile(url)
	if err != nil {
		finalerr = err
		return
	}
	fileName := getFileName(url)
	resources.StoreEmoji(fileName, name, url, data)

	imgImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		finalerr = err
		return
	}
	AllEmojis = append(AllEmojis, Emoji{name, url, []string{}, backend.ImageToRgba(imgImage)})
	e = AllEmojis[len(AllEmojis)-1]
	return
}
