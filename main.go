package main

//our imports
import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

//this function just compiles the shader dont read to much into it lol
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csrc, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csrc, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var (
	triangle = []float32{
		0, 0.5, 0, //top
		-0.5, -0.5, 0, //left
		0.5, -0.5, 0, //right
	}
)

//window constants and shader stuff
const (
	width  = 1200
	height = 800

	vertexShaderSource = `
	#version 410
	in vec3 vp;
	void main(){
		gl_Position = vec4(vp, 1.0);

	}
	` + "\x00"

	fragmentShaderSource = `
	#version 410
	out vec4 frag_colour;
	void main(){
		frag_colour = vec4(0,1,1,1);
	}
	` + "\x00"
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

	vao := buildVAO(triangle)

	//our main loop
	for !window.ShouldClose() {
		draw(window, program, vao)
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

	vs, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fs, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	//create the program,attach shaders and link it
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, fs)
	gl.LinkProgram(prog)

	//return the program
	return prog
}

//draw function
func draw(window *glfw.Window, program uint32, vao uint32) {
	//clear the screen each frame
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//tell opengl which program to use
	gl.UseProgram(program)

	//bind our vao
	gl.BindVertexArray(vao)
	//draw the triangles
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	//handle events
	glfw.PollEvents()
	//use double buffer swapping
	window.SwapBuffers()
}

//create a vao and return it
func buildVAO(points []float32) uint32 {
	//create the vbo to bind to our vao
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	//create a vao and bind it to our vbo
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexArrayAttrib(vao, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}
