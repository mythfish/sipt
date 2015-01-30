package rules

import (
	"github.com/mythfish/sipt"
)

func NewSIPInfo200Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* INFO", "SIP/2.0 200 OK", "SIP/2.0 400 Bad Request")
}
