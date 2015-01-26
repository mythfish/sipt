package sipt

import (
	"github.com/golang/glog"
	"regexp"
)

type Replacer struct {
	reg  *regexp.Regexp
	repl []byte
}

func NewReplacer(match, replace string) *Replacer {
	if match == "" {
		return nil
	}

	reg, err := regexp.Compile(match)
	if err != nil {
		glog.Warningf("Invalid replace regex: %s", err)
		return nil
	}

	repl := []byte(replace)

	//glog.Infof("Replacing %s with %s", reg.String(), repl)

	return &Replacer{
		reg:  reg,
		repl: repl,
	}
}

func (r *Replacer) replace(input []byte) []byte {
	return r.reg.ReplaceAll(input, r.repl)
}
