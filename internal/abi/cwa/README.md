

# cwa
`import "github.com/Xe/olin/internal/abi/cwa"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package cwa contains the ABI for CommonWA[1] applications.

[1]: <a href="https://github.com/CommonWA/cwa-spec">https://github.com/CommonWA/cwa-spec</a>




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func ErrorCode(err error) int](#ErrorCode)
* [type Error](#Error)
  * [func (e Error) Error() string](#Error.Error)
  * [func (i Error) String() string](#Error.String)
* [type Process](#Process)
  * [func NewProcess(name string, argv []string, env map[string]string) *Process](#NewProcess)
  * [func (p *Process) ArgAt(i int32, outPtr, outLen uint32) (int32, error)](#Process.ArgAt)
  * [func (p *Process) ArgLen() int32](#Process.ArgLen)
  * [func (p *Process) EnvGet(keyPtr, keyLen, valPtr, valLen uint32) (int32, error)](#Process.EnvGet)
  * [func (p *Process) Files() []abi.File](#Process.Files)
  * [func (p *Process) IOGetStderr() int32](#Process.IOGetStderr)
  * [func (p *Process) IOGetStdin() int32](#Process.IOGetStdin)
  * [func (p *Process) IOGetStdout() int32](#Process.IOGetStdout)
  * [func (p *Process) LogWrite(level int32, msgPtr, msgLen uint32)](#Process.LogWrite)
  * [func (p *Process) Name() string](#Process.Name)
  * [func (Process) Open(abi.File)](#Process.Open)
  * [func (p *Process) RandI32() int32](#Process.RandI32)
  * [func (p *Process) RandI64() int64](#Process.RandI64)
  * [func (p *Process) ResolveFunc(module, field string) exec.FunctionImport](#Process.ResolveFunc)
  * [func (p *Process) ResolveGlobal(module, field string) int64](#Process.ResolveGlobal)
  * [func (p *Process) ResourceClose(fid int32) error](#Process.ResourceClose)
  * [func (p *Process) ResourceFlush(fid int32) error](#Process.ResourceFlush)
  * [func (p *Process) ResourceOpen(urlPtr, urlLen uint32) (int32, error)](#Process.ResourceOpen)
  * [func (p *Process) ResourceRead(fid int32, dataPtr, dataLen uint32) (int32, error)](#Process.ResourceRead)
  * [func (p *Process) ResourceWrite(fid int32, dataPtr, dataLen uint32) (int32, error)](#Process.ResourceWrite)
  * [func (p *Process) RuntimeName(namePtr, nameLen uint32) int32](#Process.RuntimeName)
  * [func (p *Process) SetVM(vm *exec.VirtualMachine)](#Process.SetVM)
  * [func (p *Process) Setenv(m map[string]string)](#Process.Setenv)
  * [func (p *Process) TimeNow() int64](#Process.TimeNow)


#### <a name="pkg-files">Package files</a>
[c.go](/src/github.com/Xe/olin/internal/abi/cwa/c.go) [core.go](/src/github.com/Xe/olin/internal/abi/cwa/core.go) [doc.go](/src/github.com/Xe/olin/internal/abi/cwa/doc.go) [env.go](/src/github.com/Xe/olin/internal/abi/cwa/env.go) [error.go](/src/github.com/Xe/olin/internal/abi/cwa/error.go) [error_string.go](/src/github.com/Xe/olin/internal/abi/cwa/error_string.go) [io.go](/src/github.com/Xe/olin/internal/abi/cwa/io.go) [logging.go](/src/github.com/Xe/olin/internal/abi/cwa/logging.go) [random.go](/src/github.com/Xe/olin/internal/abi/cwa/random.go) [resource.go](/src/github.com/Xe/olin/internal/abi/cwa/resource.go) [runtime.go](/src/github.com/Xe/olin/internal/abi/cwa/runtime.go) [startup.go](/src/github.com/Xe/olin/internal/abi/cwa/startup.go) [time.go](/src/github.com/Xe/olin/internal/abi/cwa/time.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    LogLevelError   = 1
    LogLevelWarning = 3
    LogLevelInfo    = 6
)
```
Log levels




## <a name="ErrorCode">func</a> [ErrorCode](/src/target/error.go?s=511:540#L24)
``` go
func ErrorCode(err error) int
```
ErrorCode extracts the code from an error.




## <a name="Error">type</a> [Error](/src/target/error.go?s=132:146#L8)
``` go
type Error int
```
Error is an individual error as defined by the CommonWA spec.


``` go
const (
    ErrNone Error = (iota * -1)
    UnknownError
    InvalidArgumentError
    PermissionDeniedError
    NotFoundError
)
```
CommonWA errors as defined by the spec at <a href="https://github.com/CommonWA/cwa-spec/blob/master/errors.md">https://github.com/CommonWA/cwa-spec/blob/master/errors.md</a>










### <a name="Error.Error">func</a> (Error) [Error](/src/target/error.go?s=366:395#L19)
``` go
func (e Error) Error() string
```



### <a name="Error.String">func</a> (Error) [String](/src/target/error_string.go?s=241:271#L11)
``` go
func (i Error) String() string
```



## <a name="Process">type</a> [Process](/src/target/core.go?s=714:944#L32)
``` go
type Process struct {
    Stdin          io.Reader
    Stdout, Stderr io.Writer
    // contains filtered or unexported fields
}

```
Process is an individual CommonWA process. It is the collection of resources
and other macguffins that the child module ends up requiring.







### <a name="NewProcess">func</a> [NewProcess](/src/target/core.go?s=224:299#L15)
``` go
func NewProcess(name string, argv []string, env map[string]string) *Process
```
NewProcess creates a new process with the given name, arguments and environment.





### <a name="Process.ArgAt">func</a> (\*Process) [ArgAt](/src/target/startup.go?s=78:148#L7)
``` go
func (p *Process) ArgAt(i int32, outPtr, outLen uint32) (int32, error)
```



### <a name="Process.ArgLen">func</a> (\*Process) [ArgLen](/src/target/startup.go?s=13:45#L3)
``` go
func (p *Process) ArgLen() int32
```



### <a name="Process.EnvGet">func</a> (\*Process) [EnvGet](/src/target/env.go?s=13:91#L3)
``` go
func (p *Process) EnvGet(keyPtr, keyLen, valPtr, valLen uint32) (int32, error)
```



### <a name="Process.Files">func</a> (\*Process) [Files](/src/target/core.go?s=1475:1511#L60)
``` go
func (p *Process) Files() []abi.File
```
Files returns the set of open files in use by this process.




### <a name="Process.IOGetStderr">func</a> (\*Process) [IOGetStderr](/src/target/io.go?s=344:381#L21)
``` go
func (p *Process) IOGetStderr() int32
```



### <a name="Process.IOGetStdin">func</a> (\*Process) [IOGetStdin](/src/target/io.go?s=83:119#L9)
``` go
func (p *Process) IOGetStdin() int32
```



### <a name="Process.IOGetStdout">func</a> (\*Process) [IOGetStdout](/src/target/io.go?s=212:249#L15)
``` go
func (p *Process) IOGetStdout() int32
```



### <a name="Process.LogWrite">func</a> (\*Process) [LogWrite](/src/target/logging.go?s=31:93#L5)
``` go
func (p *Process) LogWrite(level int32, msgPtr, msgLen uint32)
```



### <a name="Process.Name">func</a> (\*Process) [Name](/src/target/core.go?s=1191:1222#L51)
``` go
func (p *Process) Name() string
```
Name returns this process' name.




### <a name="Process.Open">func</a> (Process) [Open](/src/target/core.go?s=1378:1407#L57)
``` go
func (Process) Open(abi.File)
```
Open does nothing




### <a name="Process.RandI32">func</a> (\*Process) [RandI32](/src/target/random.go?s=38:71#L7)
``` go
func (p *Process) RandI32() int32
```



### <a name="Process.RandI64">func</a> (\*Process) [RandI64](/src/target/random.go?s=98:131#L11)
``` go
func (p *Process) RandI64() int64
```



### <a name="Process.ResolveFunc">func</a> (\*Process) [ResolveFunc](/src/target/core.go?s=1803:1874#L74)
``` go
func (p *Process) ResolveFunc(module, field string) exec.FunctionImport
```
ResolveFunc resolves the CommonWA ABI and importable functions.




### <a name="Process.ResolveGlobal">func</a> (\*Process) [ResolveGlobal](/src/target/core.go?s=1662:1721#L71)
``` go
func (p *Process) ResolveGlobal(module, field string) int64
```
ResolveGlobal does nothing, currently.




### <a name="Process.ResourceClose">func</a> (\*Process) [ResourceClose](/src/target/resource.go?s=1972:2020#L91)
``` go
func (p *Process) ResourceClose(fid int32) error
```



### <a name="Process.ResourceFlush">func</a> (\*Process) [ResourceFlush](/src/target/resource.go?s=2262:2310#L108)
``` go
func (p *Process) ResourceFlush(fid int32) error
```



### <a name="Process.ResourceOpen">func</a> (\*Process) [ResourceOpen](/src/target/resource.go?s=142:210#L13)
``` go
func (p *Process) ResourceOpen(urlPtr, urlLen uint32) (int32, error)
```



### <a name="Process.ResourceRead">func</a> (\*Process) [ResourceRead](/src/target/resource.go?s=1471:1552#L69)
``` go
func (p *Process) ResourceRead(fid int32, dataPtr, dataLen uint32) (int32, error)
```



### <a name="Process.ResourceWrite">func</a> (\*Process) [ResourceWrite](/src/target/resource.go?s=1025:1107#L50)
``` go
func (p *Process) ResourceWrite(fid int32, dataPtr, dataLen uint32) (int32, error)
```



### <a name="Process.RuntimeName">func</a> (\*Process) [RuntimeName](/src/target/runtime.go?s=217:277#L17)
``` go
func (p *Process) RuntimeName(namePtr, nameLen uint32) int32
```



### <a name="Process.SetVM">func</a> (\*Process) [SetVM](/src/target/core.go?s=1293:1341#L54)
``` go
func (p *Process) SetVM(vm *exec.VirtualMachine)
```
SetVM sets the VM associated with this process.




### <a name="Process.Setenv">func</a> (\*Process) [Setenv](/src/target/core.go?s=1094:1139#L48)
``` go
func (p *Process) Setenv(m map[string]string)
```
Setenv updates a process' environment. This does not inform the process. It
is up to the running process to detect these values have changed.




### <a name="Process.TimeNow">func</a> (\*Process) [TimeNow](/src/target/time.go?s=28:61#L5)
``` go
func (p *Process) TimeNow() int64
```







- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
