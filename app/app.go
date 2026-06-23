package app

import (
	emojiManager "EmojiPicker/app/emoji"
	"fmt"

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

func Initialize() {
	clipboard.Init()
	emojiManager.Init()
}

func AfterCreateContext() {
	es := emojiManager.GetAllEMojis()
	for _, e := range es {
		var emoji emojiData
		emoji.name = e.Name
		emoji.keywords = e.Keywords
		emoji.url = e.Url
		emoji.texture = backend.NewTextureFromRgba(e.Img)
		emojis = append(emojis, emoji)
	}
}

func BeforeDestroyContext() {

}

var dockID im.ID

func Loop() {
	dockID = im.IDStr("My Dockspace")
	im.DockSpaceOverViewportV(dockID, im.MainViewport(), im.DockNodeFlagsNone, im.NewEmptyWindowClass())

	ShowPictureLoadingDemo()
}

var imageLinkInput string
var emojiNameInput string

func ShowPictureLoadingDemo() {
	im.Begin("Image")

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

				var emoji emojiData
				emoji.name = e.Name
				emoji.keywords = e.Keywords
				emoji.url = e.Url
				emoji.texture = backend.NewTextureFromRgba(e.Img)
				emojis = append(emojis, emoji)

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

	for _, emoji := range emojis {
		if im.ImageButton(emoji.name, emoji.texture.ID, im.NewVec2(48, 48)) {
			str := "[" + emoji.name + "](" + emoji.url + "?size=48&animated=true&lossless=true" + ")"
			clipboard.Write(clipboard.FmtText, []byte(str))
		}
		im.SameLine()
	}

	im.End()
}
