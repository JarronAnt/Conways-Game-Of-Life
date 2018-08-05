package main

//our imports
import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"

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
	square = []float32{
		//triangle 1
		-0.5, 0.5, 0, //top
		-0.5, -0.5, 0, //left
		0.5, -0.5, 0, //right

		//triangle 2
		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)

//definiton of a cell
type cell struct {
	//drawable holds the vao
	drawable uint32

	alive          bool
	aliveNextFrame bool

	//x,y pos
	x int
	y int
}

//window constants and shader stuff
const (
	width  = 1200
	height = 1200

	//number of rows and cols
	rows = 50
	cols = 50

	threshold = 0.15

	FPS = 10

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

	//vao := buildVAO(square)
	cells := makeCells()

	//our main loop
	for !window.ShouldClose() {
		//get curr time
		t := time.Now()

		//every frame check the state of the board
		//and modify things as needed
		for x := range cells {
			for _, c := range cells[x] {
				c.checkState(cells)
			}
		}
		//the the board each frame
		draw(window, program, cells)

		//sleep for a certain amount of time
		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
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
func draw(window *glfw.Window, program uint32, cells [][]*cell) {
	//clear the screen each frame
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//tell opengl which program to use
	gl.UseProgram(program)

	//loop through matrix of cells
	for x := range cells {
		//loop through each indivdual cell row
		for _, c := range cells[x] {
			//call the draw function on each cell
			c.draw()
		}
	}

	//cells[2][3].draw()
	//handle events
	glfw.PollEvents()
	//use double buffer swapping
	window.SwapBuffers()
}

//putting  (c *cell) before the function name basically says
//this draw function is only callable through  a *cell
//think of it like a function inside a class
func (c *cell) draw() {
	//dont draw a dead cell
	if !c.alive {
		return
	}

	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
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

func makeCells() [][]*cell {
	//seed a random number generator
	rand.Seed(time.Now().UnixNano())
	//allocate and inilatalzie a list called cells
	cells := make([][]*cell, rows, rows)

	//nested for loop creates a row of cells and appends it to the
	//cells list
	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			c := newCell(x, y)

			c.alive = rand.Float64() < threshold
			c.aliveNextFrame = c.alive
			cells[x] = append(cells[x], c)
		}
	}

	//return the cells list
	return cells
}

func newCell(x, y int) *cell {
	//copy our square definiton so we can mess with it and not
	//affect every other cell thats using square
	points := make([]float32, len(square), len(square))
	copy(points, square)

	//itterate of the points
	for i := 0; i < len(points); i++ {
		var pos float32
		var size float32

		//check if we are a x or y
		switch i % 3 {
		//set the pos to be a combo of 0.5,0,-0.5
		//set the scale of each cell by making it 10% of the col
		case 0:
			size = 1.0 / float32(cols)
			pos = float32(x) * size
		case 1:
			size = 1.0 / float32(rows)
			pos = float32(y) * size
		default:
			continue
		}

		//normalize the position
		if points[i] < 0 {
			points[i] = (pos * 2) - 1
		} else {
			points[i] = ((pos + size) * 2) - 1
		}
	}

	//return a refrence to a cell
	return &cell{
		drawable: buildVAO(points),

		x: x,
		y: y,
	}
}

//chec state takes in a gameboard and modifies the state of a isngle cell
func (c *cell) checkState(cells [][]*cell) {
	c.alive = c.aliveNextFrame
	c.aliveNextFrame = c.alive

	liveCount := c.liveNeighbours(cells)
	if c.alive {
		//if less than two neighbours are alive our cell dies
		//if more than 3 neighbours are alive cell dies
		if liveCount < 2 || liveCount > 3 {
			c.aliveNextFrame = false
		}

		//if 2 or 3 cells are alive we stay alive
		if liveCount == 2 || liveCount == 3 {
			c.aliveNextFrame = true
		}
	} else {
		//if a cell is dead but has exactly 3 neighbours its alive
		if liveCount == 3 {
			c.aliveNextFrame = true
		}
	}
}

func (c *cell) liveNeighbours(cells [][]*cell) int {
	var liveCount int
	check := func(x, y int) {
		//if we are at an edge check the otherside of the board
		if x == len(cells) {
			x = 0
		} else if x == -1 {
			x = len(cells) - 1
		}

		if y == len(cells[x]) {
			y = 0
		} else if y == -1 {
			y = len(cells[x]) - 1
		}

		//increase the liveCount if we are alive
		if cells[x][y].alive {
			liveCount++
		}
	}

	check(c.x-1, c.y)   //check the left
	check(c.x+1, c.y)   // check the right
	check(c.x, c.y+1)   //check up
	check(c.x, c.y-1)   //check down
	check(c.x-1, c.y+1) //top left
	check(c.x+1, c.y+1) //top right
	check(c.x-1, c.y-1) //bottom left
	check(c.x+1, c.y-1) //bottom right

	return liveCount
}
