package generic_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/vasilevp/aboriginal/generic"
	"golang.org/x/tools/imports"
)

func TestGeneration(t *testing.T) {
	// log.SetOutput(ioutil.Discard)

	src := `package main

import (
	"fmt"
	"log" // this import should not be copied to .gen.go
)

type T struct{} // dummy type for use in templates

type TPLᐸTᐳ struct {
	Value T "tags are preserved" // comments are not :(
}

func (t *TPLᐸTᐳ) SayHello() {
	fmt.Printf("Hello from %T\n", t)
}

func main() {
	v := TPLᐸintᐳ{} // type TPLᐸintᐳ should be generated
	fmt.Printf("%T.Value is of type %T\n", v, v.Value)
	v.Value = int(0) // will not compile if the type is wrong
	v.SayHello()

	b := TPLᐸintᐳ{} // type TPLᐸintᐳ already exists, a redefinition should not happen
	fmt.Printf("%T.Value is of type %T\n", b, b.Value)
	b.Value = int(0) // will not compile if the type is wrong
	b.SayHello()

	s := TPLᐸfloat32ᐳ{} // type TPLᐸfloat32ᐳ is different from TPLᐸintᐳ, should be generated
	fmt.Printf("%T.Value is of type %T\n", s, s.Value)
	s.Value = float32(0) // will not compile if the type is wrong
	s.SayHello()

	log.Println("")
}
`

	intermediate := &bytes.Buffer{}

	err := generic.Process(src, intermediate, "src")
	if err != nil {
		t.Error(err)
	}

	contents := intermediate.Bytes()

	contents, err = imports.Process("/tmp/out.gen.go", contents, nil)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile("/tmp/out.go", []byte(src), 0755)
	if err != nil {
		t.Error(err)
	}
	// defer os.Remove("/tmp/out.go")

	err = ioutil.WriteFile("/tmp/out.gen.go", contents, 0755)
	if err != nil {
		t.Error(err)
	}
	// defer os.Remove("/tmp/out.gen.go")

	r, err := ioutil.ReadFile("/tmp/out.gen.go")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("INPUT:\n%s\n======\nOUTPUT:\n%s\n", src, r)

	cmd := exec.Command("go", "run", "/tmp/out.go", "/tmp/out.gen.go")
	result, err := cmd.Output()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("=====\nGO RUN RESULT:\n%s", string(result))
}
