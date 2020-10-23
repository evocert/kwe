package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/evocert/kwe/ecma/es6"
	"github.com/evocert/nodejs/console"
	"github.com/evocert/nodejs/require"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var timelimit = flag.Int("timelimit", 0, "max time to run (in seconds)")

func readSource(filename string) ([]byte, error) {
	if filename == "" || filename == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filename)
}

func load(vm *es6.Runtime, call es6.FunctionCall) es6.Value {
	p := call.Argument(0).String()
	b, err := readSource(p)
	if err != nil {
		panic(vm.ToValue(fmt.Sprintf("Could not read %s: %v", p, err)))
	}
	v, err := vm.RunScript(p, string(b))
	if err != nil {
		panic(err)
	}
	return v
}

func newRandSource() es6.RandSource {
	var seed int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
		panic(fmt.Errorf("Could not read random bytes: %v", err))
	}
	return rand.New(rand.NewSource(seed)).Float64
}

func run() error {
	filename := flag.Arg(0)
	src, err := readSource(filename)
	if err != nil {
		return err
	}

	if filename == "" || filename == "-" {
		filename = "<stdin>"
	}

	vm := es6.New()
	vm.SetRandSource(newRandSource())

	new(require.Registry).Enable(vm)
	console.Enable(vm)

	vm.Set("load", func(call es6.FunctionCall) es6.Value {
		return load(vm, call)
	})

	vm.Set("readFile", func(name string) (string, error) {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return "", err
		}
		return string(b), nil
	})

	if *timelimit > 0 {
		time.AfterFunc(time.Duration(*timelimit)*time.Second, func() {
			vm.Interrupt("timeout")
		})
	}

	//log.Println("Compiling...")
	prg, err := es6.Compile(filename, string(src), false)
	if err != nil {
		return err
	}
	//log.Println("Running...")
	_, err = vm.RunProgram(prg)
	//log.Println("Finished.")
	return err
}

func main() {
	defer func() {
		if x := recover(); x != nil {
			debug.Stack()
			panic(x)
		}
	}()
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if err := run(); err != nil {
		//fmt.Printf("err type: %T\n", err)
		switch err := err.(type) {
		case *es6.Exception:
			fmt.Println(err.String())
		case *es6.InterruptedError:
			fmt.Println(err.String())
		default:
			fmt.Println(err)
		}
		os.Exit(64)
	}
}
