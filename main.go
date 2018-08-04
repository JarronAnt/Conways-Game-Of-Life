package main

//our imports
import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

//window constatns
const (
	width  = 1200
	height = 800
)

func main() {
	//make sure we run on the same thread
	runtime.LockOSThread()

	//init GLFW
	window := initGlfw()
	//this runs after surrounding functions return
	defer glfw.Terminate()

	//init opengl
	program := initOpenGL()

	//our main loop
	for !window.ShouldClose() {
		draw(window, program)
	}
}

//init glfw and return a window
func initGlfw() *glfw.Window {
	//init and check for error
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

//init opengl and return an opengl program
func initOpenGL() uint32 {
	//init and check for error
	if err := gl.Init(); err != nil {
		panic(err)
	}

	//get the opengl version and print it out
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL Version", version)

	//create the program and link it
	prog := gl.CreateProgram()
	gl.LinkProgram(prog)

	//return the program
	return prog
}

//draw function
func draw(window *glfw.Window, program uint32) {
	//clear the screen each frame
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//tell opengl which program to use
	gl.UseProgram(program)

	//handle events
	glfw.PollEvents()
	//use double buffer swapping
	window.SwapBuffers()
}
