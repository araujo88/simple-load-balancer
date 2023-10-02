package lb

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	URL           *url.URL      `json:"url"`
	ActiveConns   int32         `json:"active_conns"`
	ResponseTime  time.Duration `json:"-"`
	ResponseMutex sync.Mutex    `json:"-"`
	Weight        int           `json:"weight"`
}

type LoadBalancer struct {
	servers       []*Server
	serverLock    sync.Mutex
	algorithmType string
}

func NewLoadBalancer(filename, algorithmType string) *LoadBalancer {
	lb := &LoadBalancer{algorithmType: algorithmType}
	lb.readJson(filename)
	return lb
}

func (lb *LoadBalancer) readJson(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	err = json.Unmarshal(byteValue, &lb.servers)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}
}

// Load balancing logic: Random
func (lb *LoadBalancer) getRandomBackend() *Server {
	index := rand.Intn(len(lb.servers))
	return lb.servers[index]
}

// Load balancing logic: Least Connections
func (lb *LoadBalancer) getLeastConnectionsBackend() *Server {
	lb.serverLock.Lock()
	defer lb.serverLock.Unlock()

	// Find the backend with the least connections
	var leastConnServer *Server
	for _, server := range lb.servers {
		if leastConnServer == nil || atomic.LoadInt32(&server.ActiveConns) < atomic.LoadInt32(&leastConnServer.ActiveConns) {
			leastConnServer = server
		}
	}

	// Increment the active connection count for this server
	atomic.AddInt32(&leastConnServer.ActiveConns, 1)

	return leastConnServer
}

// Load balancing logic: Weighted Random
func (lb *LoadBalancer) getWeightedRandomBackend() *Server {
	totalWeight := 0
	for _, server := range lb.servers {
		totalWeight += server.Weight
	}

	randomValue := rand.Intn(totalWeight)

	for _, server := range lb.servers {
		randomValue -= server.Weight
		if randomValue < 0 {
			return server
		}
	}

	return lb.servers[len(lb.servers)-1] // fallback, should not happen
}

// Load balancing logic: Weighted Random Algorithm combined with Least Connections
func (lb *LoadBalancer) getDynamicWeightedRandomBackend() *Server {
	// Calculate the total dynamic weight
	totalWeight := int32(0)
	for _, server := range lb.servers {
		// Assume the weight is inversely proportional to the number of active connections
		// Add 1 to avoid division by zero
		weight := 1 / (atomic.LoadInt32(&server.ActiveConns) + 1)
		totalWeight += int32(weight)
	}

	// Select a backend server based on dynamic weight
	randomValue := int32(rand.Intn(int(totalWeight)))
	for _, server := range lb.servers {
		weight := 1 / (atomic.LoadInt32(&server.ActiveConns) + 1)
		randomValue -= int32(weight)
		if randomValue < 0 {
			atomic.AddInt32(&server.ActiveConns, 1) // Increment active connections
			return server
		}
	}

	return nil // fallback, should not happen
}

// Load balancing logic: Least Response Time
func (lb *LoadBalancer) getLeastResponseTimeBackend() *Server {
	lb.serverLock.Lock()
	defer lb.serverLock.Unlock()

	var minResponseTimeServer *Server
	for _, server := range lb.servers {
		server.ResponseMutex.Lock()
		if minResponseTimeServer == nil || server.ResponseTime < minResponseTimeServer.ResponseTime {
			minResponseTimeServer = server
		}
		server.ResponseMutex.Unlock()
	}

	// Increment the active connection count for this server
	atomic.AddInt32(&minResponseTimeServer.ActiveConns, 1)

	return minResponseTimeServer
}

// Decrement connection count for a server
func releaseBackend(server *Server) {
	atomic.AddInt32(&server.ActiveConns, -1)
}

// HTTP handler
func ProxyHandler(lb *LoadBalancer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var backend *Server

		switch lb.algorithmType {
		case "random":
			backend = lb.getRandomBackend()
		case "leastconn":
			backend = lb.getLeastConnectionsBackend()
		case "weightrand":
			backend = lb.getWeightedRandomBackend()
		case "dynamic":
			backend = lb.getDynamicWeightedRandomBackend()
		case "leasttime":
			backend = lb.getLeastResponseTimeBackend()
		default:
			log.Fatalf("Invalid algorithm type: %s", lb.algorithmType)
		}

		start := time.Now()
		proxy := httputil.NewSingleHostReverseProxy(backend.URL)
		proxy.ServeHTTP(w, r)
		duration := time.Since(start)

		backend.ResponseMutex.Lock()
		backend.ResponseTime = duration
		backend.ResponseMutex.Unlock()

		releaseBackend(backend)
	}
}
