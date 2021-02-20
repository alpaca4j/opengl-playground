package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"strings"
)

const winWidth = 640
const winHeight = 480

type ProgramID uint32
type ShaderID uint32
type BufferID uint32

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	defer sdl.Quit()
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	window, err := sdl.CreateWindow("Title", 200, 200, winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	//GetString return array of characters, GoStr truns it to a string the go can cread
	printVersion()

	//shaderProgram,_ := createProgram("resources/shaders/vertexShader.vs","resources/shaders/fragmentShader.fs")
	shaderProgram, err := NewShader("resources/shaders/vertexShader.vs", "resources/shaders/fragmentShader.fs")
	//defer gl.DeleteProgram(uint32(shaderProgram.id))
	vertices := []float32{
		0.5, 0.5, 0.0, 1.0, 1.0, //2
		0.5, -0.5, 0.0, 1.0, 0.0, //1
		-0.5, -0.5, 0.0, 0.0, 0.0, //0
		-0.5, 0.5, 0.0, 0.0, 1.0, //3

	}

	indecies := []uint32{
		0, 1, 3, //triangle 1
		1, 2, 3, //triangle 2
	}

	GenBindBuffer(gl.ARRAY_BUFFER) //VBO
	VAO := GenBindVertexArray()
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indecies)*4, gl.Ptr(indecies), gl.STATIC_DRAW)

	//Telling the VAO the format of the data
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0) //Unbind

	var r float32 = 0.0
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		shaderProgram.setFloat("r", r)

		gl.ClearColor(0.0, 1.0, 1.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		//gl.UseProgram(uint32(shaderProgram))
		shaderProgram.Use()
		gl.BindVertexArray(uint32(VAO))
		//gl.DrawArrays(gl.TRIANGLES, 0, 3) //Without index buffer,
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0)) //- with index buffer

		window.GLSwap()
		shaderProgram.CheckShadersForChanges()
		r = r + .001
	}

}

func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func GenBindBuffer(target uint32) BufferID {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	return BufferID(buffer)
}

func createProgram(vertexShaderPath string, fragmentShaderPath string) (ProgramID, error) {
	shaderProgram := gl.CreateProgram()

	vertexShader, err := createShader(gl.VERTEX_SHADER, vertexShaderPath)
	if err != nil {
		return 0, err
	}
	gl.AttachShader(shaderProgram, uint32(vertexShader))
	fragmentShader, err := createShader(gl.FRAGMENT_SHADER, fragmentShaderPath)
	if err != nil {
		return 0, err
	}
	gl.AttachShader(shaderProgram, uint32(fragmentShader))

	gl.LinkProgram(shaderProgram)
	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link programm %s", log)
	}

	gl.DeleteShader(uint32(vertexShader))
	gl.DeleteShader(uint32(fragmentShader))

	return ProgramID(shaderProgram), nil
}

func createShader(shaderType uint32, shaderFilePath string) (ShaderID, error) {

	shaderFileBytes, _ := ioutil.ReadFile(shaderFilePath)
	shaderFileBytes = append(shaderFileBytes, 0)
	shaderSource := string(shaderFileBytes)

	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("Failed to compile shader: %s", log)
	}
	return ShaderID(shader), nil

}

func printVersion() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(version)
}
