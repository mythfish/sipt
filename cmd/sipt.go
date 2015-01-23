package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/mythfish/sipt"
)

var localAddr = flag.String("l", ":9999", "local address")
var remoteAddr = flag.String("r", "localhost:80", "remote address")
var verbose = flag.Bool("v", false, "display server actions")
var veryverbose = flag.Bool("vv", false, "display server actions and all tcp data")
var nagles = flag.Bool("n", false, "disable nagles algorithm")
var match = flag.String("match", "", "match regex (in the form 'regex')")
var replace = flag.String("replace", "", "replace regex (in the form 'regex~replacer')")

func main() {
	flag.Parse()
	glog.Info("Proxying from %v to %v\n", *localAddr, *remoteAddr)

	pm := sipt.NewProxyManager(localAddr, remoteAddr, *nagles)
	go pm.Run()

}
