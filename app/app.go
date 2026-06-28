package app

import (
	emojiManager "EmojiPicker/app/emoji"
	"fmt"
	"slices"
	"strings"

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

func Initialize() {
	clipboard.Init()
	emojiManager.Init()
}

func AfterCreateContext() {
	es := emojiManager.GetAllEMojis()
	for _, e := range es {
		addEmoji(e)
	}
	sortEmojis()
}

func BeforeDestroyContext() {

}

var dockID im.ID

func Loop() {
	dockID = im.IDStr("My Dockspace")
	im.DockSpaceOverViewportV(dockID, im.MainViewport(), im.DockNodeFlagsNone, im.NewEmptyWindowClass())

	ShowEmojis()
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
