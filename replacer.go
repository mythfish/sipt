package sipt

import (
	"github.com/golang/glog"
	"regexp"
)

type Matcher struct {
	reg     *regexp.Regexp
	matchid int
}

type Replacer struct {
	reg  *regexp.Regexp
	repl []byte
}

func NewMatcher(match string) *Matcher {
	if match == "" {
		return nil
	}
	reg, err := regexp.Compile(match)
	if err != nil {
		glog.Warning("Invalid match regex: %s", err)
		return nil
	}
	glog.Info("Matching %s", reg.String())

	return &Matcher{
		reg:     reg,
		matchid: 0,
	}
}

func (m *Matcher) match(input []byte) bool {

	ms := m.reg.FindAll(input, -1)
	for _, mc := range ms {
		m.matchid++
		glog.Info("Match #%d: %s", m.matchid, string(mc))
	}
	if m.matchid == 0 {
		return false
	} else {
		return true
	}
}

func NewReplacer(match, replace string) *Replacer {
	if match == "" {
		return nil
	}

	reg, err := regexp.Compile(match)
	if err != nil {
		glog.Warning("Invalid replace regex: %s", err)
		return nil
	}

	repl := []byte(replace)

	glog.Info("Replacing %s with %s", reg.String(), repl)

	return &Replacer{
		reg:  reg,
		repl: repl,
	}
}

func (r *Replacer) replace(input []byte) []byte {
	return r.reg.ReplaceAll(input, r.repl)
}
