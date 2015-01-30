package rules

import (
	"github.com/mythfish/sipt"
)

func NewSIPInvite407Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* INVITE", "SIP/2.0 407 Proxy Authentication Required", "SIP/2.0 400 Bad Request")
}

func NewSIPInvite100Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* INVITE", "SIP/2.0 100 Trying", "SIP/2.0 502 Bad Gateway")
}

func NewSIPInvite183Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* INVITE", "SIP/2.0 183 Session Progress", "SIP/2.0 486 Busy Here")
}

func NewSIPInvite200Rule() *sipt.Rule {
	return sipt.NewRule("CSeq: \\d* INVITE", "SIP/2.0 200 OK", "SIP/2.0 503 Service Unavailable")
}
