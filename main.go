package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Server struct {
	URL         *url.URL
	ActiveConns int32
}

// List of backend servers
var servers = []*Server{
	{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:8081",
		},
		ActiveConns: 0,
	},
	{
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:8082",
		},
		ActiveConns: 0,
	},
}

var serverLock sync.Mutex

// Load balancing logic: Random
func getRandomBackend() *Server {
	index := rand.Intn(len(servers))
	return servers[index]
}

// Load balancing logic: Least Connections
func getLeastConnectionsBackend() *Server {
	serverLock.Lock()
	defer serverLock.Unlock()

	// Find the backend with the least connections
	var leastConnServer *Server
	for _, server := range servers {
		if leastConnServer == nil || atomic.LoadInt32(&server.ActiveConns) < atomic.LoadInt32(&leastConnServer.ActiveConns) {
			leastConnServer = server
		}
	}

	// Increment the active connection count for this server
	atomic.AddInt32(&leastConnServer.ActiveConns, 1)

	return leastConnServer
}

// Assume weights are stored in an array `weights`
// such that weights[i] corresponds to servers[i]
var weights = []int{3, 1} // example weights

// Load balancing logic: Weighted Random
func getWeightedRandomBackend() *Server {
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	randomValue := rand.Intn(totalWeight)

	for i, weight := range weights {
		randomValue -= weight
		if randomValue < 0 {
			return servers[i]
		}
	}

	return servers[len(servers)-1] // fallback, should not happen
}

// Load balancing logic: Weighted Random Algorithm combined with Least Connections
func getDynamicWeightedRandomBackend() *Server {
	// Calculate the total dynamic weight
	totalWeight := int32(0)
	for _, server := range servers {
		// Assume the weight is inversely proportional to the number of active connections
		// Add 1 to avoid division by zero
		weight := 1 / (atomic.LoadInt32(&server.ActiveConns) + 1)
		totalWeight += int32(weight)
	}

	// Select a backend server based on dynamic weight
	randomValue := int32(rand.Intn(int(totalWeight)))
	for _, server := range servers {
		weight := 1 / (atomic.LoadInt32(&server.ActiveConns) + 1)
		randomValue -= int32(weight)
		if randomValue < 0 {
			atomic.AddInt32(&server.ActiveConns, 1) // Increment active connections
			return server
		}
	}

	return nil // fallback, should not happen
}

// Decrement connection count for a server
func releaseBackend(server *Server) {
	atomic.AddInt32(&server.ActiveConns, -1)
}

// HTTP handler
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	backend := getDynamicWeightedRandomBackend()
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)
	proxy.ServeHTTP(w, r)
	releaseBackend(backend)
}

func main() {
	http.HandleFunc("/", proxyHandler)

	fmt.Println("Load Balancer running on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
