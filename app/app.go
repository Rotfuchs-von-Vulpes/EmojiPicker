package app

import (
	emojiManager "EmojiPicker/app/emoji"
	"EmojiPicker/app/emoji/resources"
	gifManager "EmojiPicker/app/gif"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/AllenDang/cimgui-go/backend"
	im "github.com/AllenDang/cimgui-go/imgui"
	"golang.design/x/clipboard"
)

type emojiData struct {
	name     string
	url      string
	keywords []string
	texture  *backend.Texture
}

var emojis []emojiData

func addEmoji(e emojiManager.Emoji) {
	var emoji emojiData
	emoji.name = e.Name
	emoji.keywords = e.Keywords
	emoji.url = e.Url
	emoji.texture = backend.NewTextureFromRgba(e.Img)
	emojis = append(emojis, emoji)
}

func sortEmojis() {
	less := func(a, b emojiData) int {
		return strings.Compare(a.name, b.name)
	}

	slices.SortFunc(emojis, less)
}

type gifData struct {
	id       string
	url      string
	keywords []string
	texture  *backend.Texture

	loopCount int
	delay     []int

	deltaTime int64
	last      int64
	idx       int
}

func (s *gifData) loop() {
	now := time.Now().UnixMilli()
	elapsed := now - s.last
	s.deltaTime += elapsed
	if s.idx >= s.loopCount-1 {
		s.idx = 0
	}
	delay := 10 * int64(s.delay[s.idx])
	if s.deltaTime > delay {
		s.deltaTime -= delay
		s.idx++
	}
	s.last = now
}

func (s *gifData) uv0() (v im.Vec2) {
	v.X = float32(s.idx) / float32(s.loopCount)
	v.Y = 0
	return
}

func (s *gifData) uv1() (v im.Vec2) {
	v.X = float32(s.idx+1) / float32(s.loopCount)
	v.Y = 1
	return
}

var gifs []*gifData

func addGif(g gifManager.Gif) {
	gif := new(gifData)
	gif.id = g.Id
	gif.keywords = g.Keywords
	gif.url = g.Url

	gif.texture = backend.NewTextureFromRgba(g.Img)

	gif.delay = g.Delay
	gif.loopCount = g.LoopCount
	gif.last = time.Now().UnixMilli()

	gifs = append(gifs, gif)
}

func sortGifs() {
	less := func(a, b emojiData) int {
		return strings.Compare(a.name, b.name)
	}

	slices.SortFunc(emojis, less)
}

const FPS = 60
const minimumRate int64 = 1000 / FPS

func Initialize() {
	clipboard.Init()
	emojiManager.Init()
	gifManager.Init()
}

func AfterCreateContext() {
	for _, e := range emojiManager.GetAllEMojis() {
		addEmoji(e)
	}
	sortEmojis()
	for _, g := range gifManager.GetAllGifs() {
		addGif(g)
	}
	im.CurrentIO().SetIniFilename(filepath.Join(resources.AppDir, "imgui.ini"))
}

func BeforeDestroyContext() {

}

var dockID im.ID

func Loop() {
	now := time.Now()

	dockID = im.IDStr("My Dockspace")
	im.DockSpaceOverViewportV(dockID, im.MainViewport(), im.DockNodeFlagsNone, im.NewEmptyWindowClass())

	ShowEmojis()
	ShowGifs()

	elapsed := time.Since(now).Milliseconds()
	if elapsed < minimumRate {
		time.Sleep(time.Duration(minimumRate-elapsed) * time.Millisecond)
	}
}

var imageLinkInput string
var emojiNameInput string

func ShowEmojis() {
	im.Begin("Emojis")

	if im.Button("Add Emoji") {
		im.OpenPopupStr("Add Emoji")
	}

	if im.BeginPopupModal("Add Emoji") {
		im.InputTextWithHint("Image Link", "", &imageLinkInput, im.InputTextFlagsNone, nil)
		im.InputTextWithHint("Emoji Name", "", &emojiNameInput, im.InputTextFlagsNone, nil)
		if im.Button("Add") {
			if e, err := emojiManager.SaveEmoji(imageLinkInput, emojiNameInput); err != nil {
				fmt.Println(err)
			} else {
				imageLinkInput = ""
				emojiNameInput = ""

				addEmoji(e)
				sortEmojis()

				im.CloseCurrentPopup()
			}
		}
		im.SameLine()
		if im.Button("Cancel") {
			imageLinkInput = ""
			emojiNameInput = ""
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	const width = 64
	for i, emoji := range emojis {
		availableSpace := im.ContentRegionAvail().X
		if im.ImageButton(emoji.name, emoji.texture.ID, im.NewVec2(48, 48)) {
			str := "[" + emoji.name + "](" + emoji.url + "?size=48&animated=true&lossless=true" + ")"
			clipboard.Write(clipboard.FmtText, []byte(str))
		}
		if i != len(emojis)-1 && availableSpace-2*width > 0 {
			im.SameLine()
		}
	}

	im.End()
}

var gifLinkInput string

func ShowGifs() {
	im.Begin("Gifs")

	if im.Button("Add Gif") {
		im.OpenPopupStr("Add Gif")
	}

	if im.BeginPopupModal("Add Gif") {
		im.InputTextWithHint("Image Link", "", &gifLinkInput, im.InputTextFlagsNone, nil)
		if im.Button("Add") {
			if e, err := gifManager.SaveGif(gifLinkInput); err != nil {
				fmt.Println(err)
			} else {
				gifLinkInput = ""

				addGif(e)
				sortGifs()

				im.CloseCurrentPopup()
			}
		}
		im.SameLine()
		if im.Button("Cancel") {
			gifLinkInput = ""
			im.CloseCurrentPopup()
		}
		im.EndPopup()
	}

	const width = 64
	for i, gif := range gifs {
		gif.loop()
		availableSpace := im.ContentRegionAvail().X
		if im.ImageButtonV(gif.id, gif.texture.ID, im.NewVec2(48, 48), gif.uv0(), gif.uv1(), im.NewVec4(0, 0, 0, 0), im.NewVec4(1, 1, 1, 1)) {
			str := gif.url
			clipboard.Write(clipboard.FmtText, []byte(str))
		}
		if im.IsItemHovered() {
			if im.BeginTooltip() {
				im.Text("poposa")
			}
			im.EndTooltip()
		}
		if i != len(gifs)-1 && availableSpace-2*width > 0 {
			im.SameLine()
		}
	}

	im.End()
}
