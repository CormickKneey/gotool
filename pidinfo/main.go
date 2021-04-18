package main

import (
	"flag"
)

var (
	pid    	= flag.String("p", "0", "set target pid")
	name    = flag.String("n", "all", "field name of the pid")
)


func main(){
	flag.Parse()

}

