package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/araujo88/simple-load-balancer/lb"
)

func main() {
	var algorithmType = flag.String("algorithm", "leastconn", "load balancer algorithm (options: random, leastconn, weightrand, dynamic, leasttime)")
	flag.Parse() // parse the command-line flags

	fmt.Println("Initializing load balancer - algorithm type: " + *algorithmType)

	loadBalancer := lb.NewLoadBalancer("servers.json", *algorithmType)

	http.HandleFunc("/", lb.ProxyHandler(loadBalancer))

	fmt.Println("Load Balancer running on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
