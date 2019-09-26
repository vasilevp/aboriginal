# Aboriginal Generics: the future is here!

_Inspired by this gem of an idea (click the image to go to the original comment):_

<a href="https://www.reddit.com/r/rust/comments/5penft/parallelizing_enjarify_in_go_and_rust/dcsgk7n/" ><img src=https://i.imgur.com/QSF9e4f.png height=300></a>

## Installation
```
go get github.com/vasilevp/aboriginal/cmd/aboriginal
```

## Example Usage

### command line
```sh
aboriginal -i ᐸinput_fileᐳ -o ᐸoutput_fileᐳ
```
You can omit both `-i` and `-o`, in which case STDIN/STDOUT will be used.

### go:generate example

Create `main.go` with the following contents:
```go
//go:generate aboriginal -i main.go -o main.gen.go

package main

import "fmt"

type T interface{} // template argument placeholder

type OptionalᐸTᐳ struct {
	Value T
	Valid bool
}

func main() {
	optInt := Optionalᐸintᐳ{
		Value: 42,
		Valid: true,
	}

	fmt.Printf("%T %+v\n", optInt, optInt)

	optFloat := Optionalᐸfloat64ᐳ{
		Value: 42.42,
		Valid: true,
	}

	fmt.Printf("%T %+v\n", optFloat, optFloat)
}
```

Then run it by calling
```sh
go generate && go run *.go
```

You should get the following output:
```go
main.Optionalᐸintᐳ {Value:42 Valid:true}
main.Optionalᐸfloat64ᐳ {Value:42.42 Valid:true}
```

## How does it work?
The algorithm is fairly simple:
1. Parse the source file
2. Remember all struct declarations of the form `XᐸTᐳ`
3. For any reference to `XᐸYᐳ`, where `Y` is an arbitrary type, generate an implementation of `XᐸTᐳ`, replacing `T` with `Y` for all member types (verbatim)
4. Additionally, if there are any known methods defined for `XᐸTᐳ`, generate implementations for those as well, replacing `XᐸTᐳ` with `XᐸYᐳ` in receiver type
5. Run `imports.Process()` (the core of `goimports`) on the resulting file to fix any unused imports (necessary since all imports are copied verbatim from the original file)

## TODO
- [x] Implement basic generic generation
- [x] Implement generic methods support
- [x] Implement proper import handling (sort of works)
- [ ] Implement generic functions

## Disclaimer
This project is a joke. Please, for the love of go, don't use in production.
