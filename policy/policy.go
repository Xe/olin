package policy

import (
	"flag"
	"regexp"

	"github.com/mkmik/stringlist"
	"within.website/confyg/flagconfyg"
)

type Policy struct {
	Allowed      []*regexp.Regexp
	Disallowed   []*regexp.Regexp
	RamPageLimit uint
	GasLimit     uint
}

func Parse(name string, data []byte) (Policy, error) {
	var result Policy
	fs := flag.NewFlagSet(name, flag.ContinueOnError)

	var allowed stringlist.Value
	var disallowed stringlist.Value
	fs.Var(&allowed, "allow", "the list of file URL's that are allowed")
	fs.Var(&disallowed, "disallow", "the list of file URL's that are disallowed")
	fs.UintVar(&result.RamPageLimit, "ram-page-limit", 128, "the ram page limit")
	fs.UintVar(&result.GasLimit, "gas-limit", 32*1024*1024, "the number of wasm instructions that can be run")

	err := flagconfyg.Parse(name, data, fs)
	if err != nil {
		return Policy{}, err
	}

	for _, allow := range []string(allowed) {
		rex, err := regexp.Compile(allow)
		if err != nil {
			return Policy{}, err
		}
		result.Allowed = append(result.Allowed, rex)
	}

	for _, disallow := range []string(disallowed) {
		rex, err := regexp.Compile(disallow)
		if err != nil {
			return Policy{}, err
		}
		result.Disallowed = append(result.Disallowed, rex)
	}

	return result, nil
}
