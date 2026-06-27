package main

import (
	"runtime"

	"EmojiPicker/app"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
)

var currentBackend backend.Backend[sdlbackend.SDLWindowFlags]

func init() {
	runtime.LockOSThread()
}

func main() {
	app.Initialize()

	currentBackend, _ = backend.CreateBackend(sdlbackend.NewSDLBackend())
	currentBackend.SetAfterCreateContextHook(app.AfterCreateContext)
	currentBackend.SetBeforeDestroyContextHook(app.BeforeDestroyContext)

	currentBackend.SetBgColor(imgui.NewVec4(0.45, 0.55, 0.6, 1.0))

	currentBackend.CreateWindow("Fake Emoji Picker", 1200, 900)

	currentBackend.SetSwapInterval(sdlbackend.SDLSwapIntervalVsync)

	currentBackend.Run(app.Loop)
}
