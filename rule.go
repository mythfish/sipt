package sipt

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
	if r.matcher.match(input) {
		return r.replacer.replace(input)
	}
	return input
}

func (r *Rule) Key() string {
	return r.matcher.String()
}
