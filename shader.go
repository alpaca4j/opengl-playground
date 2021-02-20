package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"os"
	"time"
)

type Shader struct {
	id               ProgramID
	vertexPath       string
	fragmentPath     string
	vertexModified   time.Time
	fragmentModified time.Time
}

func NewShader(vertexPath string, fragmentPath string) (*Shader, error) {
	programId, err := createProgram(vertexPath, fragmentPath)
	if err != nil {
		return nil, err
	}

	result := &Shader{
		id:               programId,
		vertexPath:       vertexPath,
		fragmentPath:     fragmentPath,
		vertexModified:   GetModifiecTime(vertexPath),
		fragmentModified: GetModifiecTime(fragmentPath),
	}

	return result, nil
}

func GetModifiecTime(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		panic(err)
	}
	return info.ModTime()
}

func (shader *Shader) Use() {
	gl.UseProgram(uint32(shader.id))
}

func (shader *Shader) setFloat(name string, f float32) {
	name_cstr := gl.Str(name + "\x00") //C style string -> array of bytes
	location := gl.GetUniformLocation(uint32(shader.id), name_cstr)
	gl.Uniform1f(location, f)
}

func (sh *Shader) CheckShadersForChanges() {

	vertexModTime := GetModifiecTime(sh.vertexPath)
	fragmentModTime := GetModifiecTime(sh.fragmentPath)

	if !sh.vertexModified.Equal(vertexModTime) ||
		!sh.fragmentModified.Equal(fragmentModTime) {
		fmt.Print("Found changes in time")
		id, err := createProgram(sh.vertexPath, sh.fragmentPath)
		if err != nil {
			fmt.Println(err)
		}
		sh.fragmentModified = fragmentModTime
		sh.vertexModified = vertexModTime
		gl.DeleteProgram(uint32(sh.id))
		sh.id = id

	}

}
