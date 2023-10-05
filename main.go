package main

import (
	"flag"
	"fmt"

	"github.com/araujo88/simple-load-balancer/lb"
	"github.com/valyala/fasthttp"
)

func main() {
	var algorithmType = flag.String("algorithm", "leastconn", "load balancer algorithm (options: random, leastconn, weightrand, dynamic, dynamic2, leasttime)")
	flag.Parse() // parse the command-line flags

	fmt.Println("Initializing load balancer - algorithm type: " + *algorithmType)

	loadBalancer := lb.NewLoadBalancer("servers.json", *algorithmType)

	// Set up the fasthttp server
	server := &fasthttp.Server{
		Handler: lb.ProxyHandler(loadBalancer),
	}

	fmt.Println("Load Balancer running on :3000")
	if err := server.ListenAndServe(":3000"); err != nil {
		panic(err)
	}
}
