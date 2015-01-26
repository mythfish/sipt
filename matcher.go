package sipt

import (
	"github.com/golang/glog"
	"regexp"
)

type Matcher struct {
	reg     *regexp.Regexp
	matchid int
}

func NewMatcher(match string) *Matcher {
	if match == "" {
		return nil
	}
	reg, err := regexp.Compile(match)
	if err != nil {
		glog.Warningf("Invalid match regex: %s", err)
		return nil
	}
	//glog.Infof("Matching %s", reg.String())

	return &Matcher{
		reg:     reg,
		matchid: 0,
	}
}

func (m *Matcher) match(input []byte) bool {

	ms := m.reg.FindAll(input, -1)
	for _, mc := range ms {
		m.matchid++
		glog.Infof("Match #%d: %s", m.matchid, string(mc))
	}
	if m.matchid == 0 {
		return false
	} else {
		return true
	}
}

func (m *Matcher) String() string {
	return m.reg.String()
}
