package sipt

import (
	"github.com/golang/glog"
)

type Rule struct {
	matcher  *Matcher
	replacer *Replacer
}

func NewRule(match, rMatch, replace string) *Rule {
	return &Rule{
		matcher:  NewMatcher(match),
		replacer: NewReplacer(rMatch, replace),
	}
}

func (r *Rule) Do(input []byte) []byte {
	glog.Infof("try MATCH key:%s", r.matcher.String())
	if r.matcher.match(input) {
		glog.Infof("MATCH key:%s", r.matcher.String())
		return r.replacer.replace(input)
	}
	return input
}

func (r *Rule) Key() string {
	return r.matcher.String()
}
