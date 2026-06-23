package main

import (
	"fmt"
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

	currentBackend.CreateWindow("Hello from cimgui-go", 1200, 900)

	currentBackend.SetDropCallback(func(p []string) {
		fmt.Printf("drop triggered: %v", p)
	})

	currentBackend.SetCloseCallback(func() {
		fmt.Println("window is closing")
	})

	currentBackend.Run(app.Loop)
}
