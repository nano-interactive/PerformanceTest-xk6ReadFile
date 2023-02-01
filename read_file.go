// Package read_file provides the xk6 Modules implementation for reading a file line by line concurrently in order
// using Javascript
package read_file

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"io"
	"os"
	"sync"
)

var mu sync.Mutex
var file *os.File
var scanner *bufio.Scanner
var funcCallable goja.Callable
var started = false

func init() {
	modules.Register("k6/x/read-file", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/read-file` module instances for each VU.
type RootModule struct{}

// ReadFile represents an instance of the ReadFile module for every VU.
type ReadFile struct {
	vu modules.VU
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ReadFile{}
)

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ReadFile{vu: vu}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (r *ReadFile) Exports() modules.Exports {
	return modules.Exports{Default: r}
}

// OpenFile is a wrapper for Go read_file.OpenFile and it must be called in init context
func (r *ReadFile) OpenFile(filePath string) {
	rt := r.vu.Runtime()

	if r.vu.State() != nil && r.vu.State().VUID != 0 {
		common.Throw(rt, errors.New("OpenFile must be called in the setup function"))
	}
	var err error
	file, err = os.Open(filePath)
	if err != nil {
		common.Throw(rt, errors.New(fmt.Sprintf("error reading file %s", filePath)))
	}
	scanner = bufio.NewScanner(file)
}

// ReadLine is a wrapper for Go read_file.ReadLine
func (r *ReadFile) ReadLine() string {
	mu.Lock()
	defer mu.Unlock()

	rt := r.vu.Runtime()

	if !started {
		if funcCallable != nil {
			r.callFileStartJsFunc()
		}
		started = true
	}

	if scanner == nil {
		common.Throw(rt, errors.New("file is not opened, use OpenFile in setup function"))
	}

	if !scanner.Scan() {
		if !r.resetFilePointer() {
			common.Throw(rt, errors.New("nothing to read, the file is empty"))
		}
		if funcCallable != nil {
			r.callFileStartJsFunc()
		}

	}

	return scanner.Text()
}

func (r *ReadFile) callFileStartJsFunc() {
	var err error
	for i := 0; i < 5; i++ {
		_, err = funcCallable(goja.Undefined())
		if err == nil {
			break
		}
	}
	if err != nil {
		common.Throw(r.vu.Runtime(), err)
	}
}

// resetFilePoint resets a file pointer to the beginning of the file
func (r *ReadFile) resetFilePointer() bool {

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		common.Throw(r.vu.Runtime(), err)
	}
	scanner = bufio.NewScanner(file)
	return scanner.Scan()
}

func (r *ReadFile) SetFileStartJsFunc(f goja.Callable) {
	if f == nil {
		common.Throw(r.vu.Runtime(), errors.New("SetFileStartJsFunc requires a function as first argument"))
	}
	funcCallable = f
}

func (*ReadFile) Close() {
	_ = file.Close()
}
