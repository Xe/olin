

# dagger
`import "github.com/Xe/olin/internal/abi/dagger"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package dagger is the first attempt at an API for webassembly modules to communicate with the outside world. It is based on the idea of files being used as the intermediate between user modules and resources.

Consider this the first draft of Dagger, everything here is subject to change. This is going to be the experimental phase.




## <a name="pkg-index">Index</a>
* [type Errno](#Errno)
  * [func (i Errno) String() string](#Errno.String)
* [type Error](#Error)
  * [func (e Error) Error() string](#Error.Error)
* [type Process](#Process)
  * [func NewProcess(name string) *Process](#NewProcess)
  * [func (p *Process) CloseFD(fd int64) int64](#Process.CloseFD)
  * [func (p Process) Files() []abi.File](#Process.Files)
  * [func (p *Process) Name() string](#Process.Name)
  * [func (p *Process) Open(f abi.File)](#Process.Open)
  * [func (p *Process) OpenFD(furl string, flags uint32) int64](#Process.OpenFD)
  * [func (p *Process) ReadFD(fd int64, buf []byte) int64](#Process.ReadFD)
  * [func (p *Process) ResolveFunc(module, field string) exec.FunctionImport](#Process.ResolveFunc)
  * [func (p *Process) ResolveGlobal(module, field string) int64](#Process.ResolveGlobal)
  * [func (p *Process) SyncFD(fd int64) int64](#Process.SyncFD)
  * [func (p *Process) WriteFD(fd int64, data []byte) int64](#Process.WriteFD)


#### <a name="pkg-files">Package files</a>
[c.go](/src/github.com/Xe/olin/internal/abi/dagger/c.go) [doc.go](/src/github.com/Xe/olin/internal/abi/dagger/doc.go) [errno_string.go](/src/github.com/Xe/olin/internal/abi/dagger/errno_string.go) [error.go](/src/github.com/Xe/olin/internal/abi/dagger/error.go) [process.go](/src/github.com/Xe/olin/internal/abi/dagger/process.go) 






## <a name="Errno">type</a> [Errno](/src/target/error.go?s=419:433#L23)
``` go
type Errno int
```
Errno is the error number for an error.


``` go
const (
    ErrorNone Errno = iota
    ErrorBadURL
    ErrorBadURLInput
    ErrorUnknownScheme
)
```
Error numbers










### <a name="Errno.String">func</a> (Errno) [String](/src/target/errno_string.go?s=220:250#L11)
``` go
func (i Errno) String() string
```



## <a name="Error">type</a> [Error](/src/target/error.go?s=85:142#L6)
``` go
type Error struct {
    Errno      Errno
    Underlying error
}

```
Error is a common error type for dagger operations.










### <a name="Error.Error">func</a> (Error) [Error](/src/target/error.go?s=144:173#L11)
``` go
func (e Error) Error() string
```



## <a name="Process">type</a> [Process](/src/target/process.go?s=283:338#L18)
``` go
type Process struct {
    // contains filtered or unexported fields
}

```
Process is a higher level wrapper around a set of files for dagger
modules.







### <a name="NewProcess">func</a> [NewProcess](/src/target/process.go?s=484:521#L29)
``` go
func NewProcess(name string) *Process
```
NewProcess creates a new process.





### <a name="Process.CloseFD">func</a> (\*Process) [CloseFD](/src/target/process.go?s=1306:1347#L64)
``` go
func (p *Process) CloseFD(fd int64) int64
```
CloseFD closes a file descriptor.




### <a name="Process.Files">func</a> (Process) [Files](/src/target/process.go?s=390:425#L24)
``` go
func (p Process) Files() []abi.File
```
Files returns the process' list of open files.




### <a name="Process.Name">func</a> (\*Process) [Name](/src/target/process.go?s=868:899#L48)
``` go
func (p *Process) Name() string
```
Name returns this process's name.




### <a name="Process.Open">func</a> (\*Process) [Open](/src/target/process.go?s=3843:3877#L170)
``` go
func (p *Process) Open(f abi.File)
```
Open makes this Process track an arbitrary extra file.




### <a name="Process.OpenFD">func</a> (\*Process) [OpenFD](/src/target/process.go?s=1023:1080#L52)
``` go
func (p *Process) OpenFD(furl string, flags uint32) int64
```
OpenFD opens a file descriptor for this Process with the given file url
string and flags integer.




### <a name="Process.ReadFD">func</a> (\*Process) [ReadFD](/src/target/process.go?s=2111:2163#L100)
``` go
func (p *Process) ReadFD(fd int64, buf []byte) int64
```
ReadFD reads data from the given FD into the byte buffer.




### <a name="Process.ResolveFunc">func</a> (\*Process) [ResolveFunc](/src/target/process.go?s=2389:2460#L111)
``` go
func (p *Process) ResolveFunc(module, field string) exec.FunctionImport
```
ResolveFunc resolves dagger's ABI and importable functions.




### <a name="Process.ResolveGlobal">func</a> (\*Process) [ResolveGlobal](/src/target/process.go?s=3942:4001#L175)
``` go
func (p *Process) ResolveGlobal(module, field string) int64
```
ResolveGlobal does nothing, currently.




### <a name="Process.SyncFD">func</a> (\*Process) [SyncFD](/src/target/process.go?s=1884:1924#L89)
``` go
func (p *Process) SyncFD(fd int64) int64
```
SyncFD runs a file's sync operation and returns -1 if it failed.




### <a name="Process.WriteFD">func</a> (\*Process) [WriteFD](/src/target/process.go?s=1590:1644#L78)
``` go
func (p *Process) WriteFD(fd int64, data []byte) int64
```
WriteFD writes the given data to a file descriptor, returning -1 if it failed
somehow.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
