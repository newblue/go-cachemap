package main

//Reads hostname's from stdin and looks them up, caching results for a period of time
//It's pretty useless, but a good example of the cachemap package

import (
	"net"
	"fmt"
	"os"
	"cachemap"
	"time"
	"bufio"
)

var cache *cachemap.Cache

type Result struct {
	Cname string
	Addrs []string
}

//Lookup a hostname, write to stderr so we can see when a value is cached and when it isn't
func Lookup(name string) (interface{}, bool) {
	os.Stderr.WriteString("Looking up " + name + "\n")
	c, a, e := net.LookupHost(name)
	if e != nil {
		return nil, false
	}
	//This goroutine will mark a hostname as stale after a period of time
	go func() {
		time.Sleep(6e10) //Only holds for 1 minute
		cache.Stale(name)
	}()
	return &Result{c, a}, true
}

func main() {
	cache = cachemap.New()
	in := bufio.NewReader(os.Stdin)
	for name, e := in.ReadBytes('\n'); e == nil; name, e = in.ReadBytes('\n') {
		if name[len(name) - 1] == '\n' { name = name[0:len(name) - 1] }

		res, ok := cache.Get(string(name), Lookup)
		if !ok {
			fmt.Println("Couldn't find " + string(name))
		} else {
			fmt.Printf("%#+v\n", res.(*Result))
		}
	}
}
