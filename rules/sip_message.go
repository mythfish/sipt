package rules

import (
	"github.com/mythfish/sipt"
)

func NewSIPMessage200Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* MESSAGE", "SIP/2.0 200 OK", "SIP/2.0 502 Bad Gateway")
}
