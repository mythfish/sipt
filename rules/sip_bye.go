package rules

import (
	"github.com/mythfish/sipt"
)

func NewSIPBye200Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* BYE", "SIP/2.0 200 OK", "SIP/2.0 481 Call/Transaction Does Not Exist")
}
