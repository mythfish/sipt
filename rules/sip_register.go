package rules

import (
	"github.com/mythfish/sipt"
)

func NewSIPRegister503Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* REGISTER", "SIP/2.0 200 OK", "SIP/2.0 503 Service Unavailable")
}
