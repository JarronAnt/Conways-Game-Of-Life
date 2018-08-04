package main

//our imports
import (
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

//window constatns
const (
	width  = 500
	height = 500
)

func main() {
	//make sure we run on the same thread
	runtime.LockOSThread()

	//init GLFW
	window := initGlfw()
	//this runs after surrounding functions return
	defer glfw.Terminate()

	//our main loop
	for !window.ShouldClose() {
		//todo
	}
}

//init glfw and return a window
func initGlfw() *glfw.Window {
	//basic error handling
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	//define glfw window properties
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	//create the window
	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)

	//error handle
	if err != nil {
		panic(err)
	}

	//create the window context
	window.MakeContextCurrent()

	//return the window
	return window

}
