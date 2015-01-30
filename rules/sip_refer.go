package rules

import (
	"github.com/mythfish/sipt"
)

func NewSIPRefer200Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* REFER", "SIP/2.0 202 Accepted", "SIP/2.0 400 Bad Request")
}
