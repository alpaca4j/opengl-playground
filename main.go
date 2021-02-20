package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	"strings"
)

const winWidth =640
const winHeight = 480

type ProgramID uint32
type ShaderID uint32
type VboId uint32
type VaoId uint32

func main(){
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	defer sdl.Quit()
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK,sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION,3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION,3)

	window ,err := sdl.CreateWindow("Title",200,200,winWidth,winHeight,sdl.WINDOW_OPENGL)
	if err !=nil { panic(err)}

	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	//GetString return array of characters, GoStr truns it to a string the go can cread
	printVersion()

	shaderProgram := createProgram("resources/shaders/vertexShader.vs","resources/shaders/fragmentShader.fs")

	defer gl.DeleteProgram(uint32(shaderProgram))
	vertices := []float32{
		-0.5, -0.5, //0
		0.5, -0.5, //1
		0.5, 0.5, //2
		//-0.5,0.5, //3

	}

	//indecies := []uint32 {
	//0,1,2,
	//2,3,0,
	//}

	GenBindBuffer(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	//var IBO uint32
	//gl.GenBuffers(1,&IBO)
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER,IBO)
	//gl.BufferData(gl.ELEMENT_ARRAY_BUFFER,len(indecies)*4,gl.Ptr(indecies),gl.STATIC_DRAW)

	VAO := GenBindVertexArray()

	gl.EnableVertexAttribArray(0)
	//Telling the VAO the format of the data
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*2, nil)

	gl.BindVertexArray(0)

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		gl.ClearColor(0.0, 1.0, 1.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(uint32(shaderProgram))
		gl.BindVertexArray(uint32(VAO))
		gl.DrawArrays(gl.TRIANGLES, 0, 3) //Without index buffer,
		//gl.DrawElements(gl.TRIANGLES,6,gl.UNSIGNED_INT,gl.PtrOffset(0))  //- with index buffer

		window.GLSwap()
	}

}

func GenBindVertexArray() VaoId {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return VaoId(VAO)
}

func GenBindBuffer(target uint32) VboId {
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(target, VBO)
	return VboId(VBO)
}

func createProgram(vertexShaderPath string,fragmentShaderPath string) ProgramID {
	shaderProgram := gl.CreateProgram()

	vertexShader := uint32(createShader(gl.VERTEX_SHADER, vertexShaderPath))
	gl.AttachShader(shaderProgram, vertexShader)
	fragmentShader := uint32(createShader(gl.FRAGMENT_SHADER, fragmentShaderPath))
	gl.AttachShader(shaderProgram, fragmentShader)

	gl.LinkProgram(shaderProgram)
	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		panic(fmt.Sprintf("failed to link programm %s", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return ProgramID(shaderProgram)
}

func createShader(shaderType uint32, shaderFilePath string) ShaderID{

	shaderFileBytes, _ := ioutil.ReadFile(shaderFilePath)
	shaderFileBytes = append(shaderFileBytes, 0)
	shaderSource := string(shaderFileBytes)

	shader := gl.CreateShader(shaderType)
	csource,free := gl.Strs(shaderSource)
	gl.ShaderSource(shader,1,csource,nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader,gl.COMPILE_STATUS,&status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader,gl.INFO_LOG_LENGTH,&logLength)
		log := strings.Repeat("\x00",int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength,nil,gl.Str(log))
		panic(fmt.Sprintf("Failed to compile shader: %s", log))
	}
	free()

	return ShaderID(shader)

}

func printVersion() {
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println(version)
}