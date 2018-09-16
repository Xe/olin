package main

import (
	"log"
	"os"
	"time"

	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/Xe/olin/internal/names"
	"github.com/perlin-network/life/exec"
)

var (
	p             *cwa.Process
	vm            *exec.VirtualMachine
	nothingModule = []byte{0, 0x61, 0x73, 0x6d, 0x01, 0, 0, 0}
)

func init() {
	var err error
	p = cwa.NewProcess(os.Args[0], os.Args[1:], map[string]string{
		"POWERED_BY": "Go and ungodly hacking",
	})

	cfg := exec.VMConfig{}
	vm, err = exec.NewVirtualMachine(nothingModule, cfg, p)
	if err != nil {
		log.Fatal(err)
	}
}

// XXX not used
func main() {}

//export env_IR_log_write
func envIRLogWrite(level int32, msgPtr, msgLen uint32) {
	p.LogWrite(level, msgPtr, msgLen)
}

//export env_IR_env_get
func envIREnvGet(keyPtr, keyLen, valPtr, valLen uint32) int64 {
	result, err := p.EnvGet(keyPtr, keyLen, valPtr, valLen)
	if err != nil {
		return int64(cwa.ErrorCode(err))
	}

	return int64(result)
}

//export env_IR_runtime_spec_major
func envRuntimeSpecMajor() int64 {
	return names.CommonWASpecMajor
}

//export env_IR_runtime_spec_minor
func envRuntimeSpecMinor() int64 {
	return names.CommonWASpecMinor
}

//export env_IR_runtime_name
func envRuntimeName(namePtr, nameLen uint32) int64 {
	return int64(p.RuntimeName(namePtr, nameLen))
}

//export env_IR_runtime_msleep
func envRuntimeMsleep(ms int32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

//export env_IR_startup_arg_len
func envStartupArgLen() int64 {
	return int64(p.ArgLen())
}

//export env_IR_startup_arg_at
func envStartupArgAt(i int32, outPtr, outLen uint32) int64 {
	result, err := p.ArgAt(i, outPtr, outLen)
	if err != nil {
		return int64(cwa.ErrorCode(err))
	}

	return int64(result)
}

//export env_IR_resource_open
func envResourceOpen(urlPtr, urlLen uint32) int64 {
	result, err := p.ResourceOpen(urlPtr, urlLen)
	if err != nil {
		return int64(cwa.ErrorCode(err))
	}

	return int64(result)
}

//export env_IR_resource_read
func envResourceRead(fid int32, dataPtr, dataLen uint32) int64 {
	result, err := p.ResourceRead(fid, dataPtr, dataLen)
	if err != nil {
		return int64(cwa.ErrorCode(err))
	}

	return int64(result)
}

//export env_IR_resource_write
func envResourceWrite(fid int32, dataPtr, dataLen uint32) int64 {
	result, err := p.ResourceWrite(fid, dataPtr, dataLen)
	if err != nil {
		return int64(cwa.ErrorCode(err))
	}

	return int64(result)
}

//export env_IR_resource_close
func envResourceClose(fid int32) int64 {
	err := p.ResourceClose(fid)
	if err != nil {
		return int64(cwa.ErrorCode(err))
	}

	return 0
}

//export env_IR_resource_flush
func envResourceFlush(fid int32) int64 {
	err := p.ResourceFlush(fid)
	if err != nil {
		return int64(cwa.ErrorCode(err))
	}

	return 0
}

//export env_IR_time_now
func envTimeNow() int64 {
	return p.TimeNow()
}

//export env_IR_io_get_stdin
func envIOGetStdin() int64 {
	return int64(p.IOGetStdin())
}

//export env_IR_io_get_stdout
func envIOGetStdout() int64 {
	return int64(p.IOGetStdout())
}

//export env_IR_io_get_stderr
func envIOGetStderr() int64 {
	return int64(p.IOGetStderr())
}

//export env_IR_random_i32
func envRandomI32() int32 {
	return p.RandI32()
}

//export env_IR_random_i64
func envRandomI64() int64 {
	return p.RandI64()
}
