package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/Xe/olin/internal/namegen"
	"github.com/Xe/olin/internal/names"
	"github.com/Xe/olin/proto/brand"
	humanize "github.com/dustin/go-humanize"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/golang/protobuf/proto"
)

var (
	enableJit    = flag.Bool("enable-jit", false, "flag runtime to enable JIT?")
	defaultPages = flag.Int("default-pages", 32, "default number of WebAssembly pages to use, default is 32 (~2MB)")
	maxPages     = flag.Int("max-pages", 48, "maximum number of WebAssembly pages that can be used, default is 48 (~3MB)")
	mainFunc     = flag.String("main-func", "cwa_main", "\"main\" entrypoint of the webassembly module provided via flags")
	kind         = flag.String("kind", "cwagi", "binary kind")
	name         = flag.String("name", "", "name of the person or group who created this software")
)

func init() {
	flag.Parse()

	if *name == "" {
		cmd := exec.Command("git", "config", "user.name")

		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		myname := string(out)
		*name = myname[:len(myname)-1]
	}

	switch *kind {
	case "cwagi", "daemon", "interactive":
	default:
		log.Fatalf("unknown -kind %s: want cwagi, daemon or interactive", *kind)
	}
}

func main() {
	// flags already parsed
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(2)
	}

	fname := flag.Arg(0)

	fin, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("loading %s...", fname)
	mod, err := wasm.DecodeModule(fin)
	fin.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", mod.Customs[0].Name)

	if sc := mod.Custom("olin-settings"); sc != nil {
		log.Fatal("custom settings already exist for this binary")
	}

	log.Println("VM Settings:")
	log.Printf("JIT Enabled: %v", *enableJit)
	log.Printf("Default ram pages: %d (%s)", *defaultPages, humanize.Bytes(uint64(65536**defaultPages)))
	log.Printf("Max ram pages: %d (%s)", *maxPages, humanize.Bytes(uint64(65536**maxPages)))
	log.Printf("Expected runtime: %s", names.CommonWARuntimeName)
	log.Printf("Author: %s", *name)

	b := &brand.Brand{
		Opts: &brand.VMOptions{
			EnableJit:    *enableJit,
			DefaultPages: int32(*defaultPages),
			MaxPages:     int32(*maxPages),
			MainFunc:     *mainFunc,
		},
		Meta: &brand.Metadata{
			ExpectedRuntime: names.CommonWARuntimeName,
			Author:          *name,
			Name:            namegen.Next(),
		},
	}
	log.Printf("Name: %v", b.Meta.Name)

	data, err := proto.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}

	sc := &wasm.SectionCustom{
		Name: "olin-settings",
		Data: data,
	}

	mod.Customs = append(mod.Customs, sc)

	fout, err := os.Create(fname + ".branded")
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	err = wasm.EncodeModule(fout, mod)
	if err != nil {
		log.Fatal(err)
	}
}
