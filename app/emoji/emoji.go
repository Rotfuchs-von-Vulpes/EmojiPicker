package emoji

import (
	"EmojiPicker/app/emoji/resources"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/webp"

	"github.com/AllenDang/cimgui-go/backend"
)

func downloadFile(url string) ([]byte, error) {
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

	Animated  bool
	Delay     []int
	LoopCount int
}

var AllEmojis []Emoji

func wrapGif(gifImage *gif.GIF) (img *image.RGBA) {
	width := gifImage.Config.Width
	height := gifImage.Config.Height
	w := len(gifImage.Image) * width
	h := height
	img = image.NewRGBA(image.Rect(0, 0, w, h))

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	var prevCanvas *image.RGBA

	for i, frame := range gifImage.Image {
		bounds := frame.Bounds()
		disposal := gifImage.Disposal[i]

		if disposal == gif.DisposalPrevious {
			prevCanvas = image.NewRGBA(canvas.Bounds())
			draw.Draw(prevCanvas, canvas.Bounds(), canvas, image.Point{}, draw.Src)
		}

		draw.Draw(canvas, bounds, frame, bounds.Min, draw.Over)

		frameSnapshot := image.NewRGBA(canvas.Bounds())
		draw.Draw(frameSnapshot, canvas.Bounds(), canvas, image.Point{}, draw.Src)

		offset := i * width
		dstBounds := image.Rect(offset, 0, offset+width, height)
		draw.Draw(img, dstBounds, frameSnapshot, img.Bounds().Min, draw.Src)

		switch disposal {
		case gif.DisposalBackground:
			draw.Draw(canvas, bounds, image.Transparent, image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			if prevCanvas != nil {
				canvas = prevCanvas
			}
		}
	}
	return
}

func Init() {
	resources.Init()
	for _, e := range resources.AllEmojisData {
		if e.Animated {
			if _, err := gif.DecodeConfig(bytes.NewReader(e.Data)); err != nil {
				fmt.Println(err)
				continue
			}
			gifImage, err := gif.DecodeAll(bytes.NewReader(e.Data))
			if err != nil {
				fmt.Println(err)
				continue
			}
			var e2 Emoji
			e2.Animated = true
			e2.Name = e.Name
			e2.Url = e.Url
			e2.Img = wrapGif(gifImage)
			e2.Delay = gifImage.Delay
			e2.LoopCount = len(gifImage.Image)
			AllEmojis = append(AllEmojis, e2)
		} else {
			imgImage, _, err := image.Decode(bytes.NewReader(e.Data))
			if err != nil {
				fmt.Println(err)
				continue
			}
			AllEmojis = append(AllEmojis, Emoji{e.Name, e.Url, []string{}, backend.ImageToRgba(imgImage), false, nil, 0})
		}
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
	pair := strings.Split(url, "?")
	baseUrl := pair[0]
	downloadUrl := url
	animated := false
	if len(pair) > 1 {
		args := strings.SplitSeq(pair[1], "&")
		var newArgs []string
		var newBase = pair[0]
		for arg := range args {
			if arg == "animated=true" {
				animated = true
				newBase = strings.Replace(pair[0], "webp", "gif", 1)
			}
			if strings.Split(arg, "=")[0] != "size" {
				newArgs = append(newArgs, arg)
			}
		}
		downloadUrl = newBase + "?" + strings.Join(newArgs, "&")
	}
	dirs := strings.Split(url, "/")
	if len(dirs) < 2 || dirs[2] != "cdn.discordapp.com" {
		finalerr = fmt.Errorf("%s não é uma URL para o discord", url)
		return
	}
	data, err := downloadFile(downloadUrl)
	if err != nil {
		finalerr = err
		return
	}
	fileName := getFileName(downloadUrl)
	resources.StoreEmoji(fileName, name, baseUrl, animated, data)

	if animated {
		if _, err := gif.DecodeConfig(bytes.NewReader(data)); err != nil {
			finalerr = err
			return
		}
		gifImage, err := gif.DecodeAll(bytes.NewReader(data))
		if err != nil {
			finalerr = err
			return
		}

		e.Animated = true
		e.Name = name
		e.Delay = gifImage.Delay
		e.LoopCount = len(gifImage.Image)
		e.Img = wrapGif(gifImage)
		e.Url = baseUrl

		AllEmojis = append(AllEmojis, e)
	} else {
		imgImage, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			finalerr = err
			return
		}
		e = Emoji{name, baseUrl, []string{}, backend.ImageToRgba(imgImage), false, nil, 0}
		AllEmojis = append(AllEmojis, e)
	}
	return
}
