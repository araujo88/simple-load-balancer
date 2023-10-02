package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getCPUUtilization() (float64, error) {
	cmd := exec.Command("./bin/cpuusage") // assuming the C program is compiled and present in the same directory.
	out, err := cmd.Output()
	if err != nil {
		return 0.0, fmt.Errorf("error executing cpuusage: %w", err)
	}

	// Assuming the output is in the format "CPU Utilization: xx.xx%"
	outputStr := strings.TrimSpace(string(out))
	splitStr := strings.Split(outputStr, ": ")
	if len(splitStr) < 2 {
		return 0.0, fmt.Errorf("unexpected output format: %s", outputStr)
	}

	utilizationStr := strings.TrimSuffix(splitStr[1], "%")
	utilization, err := strconv.ParseFloat(utilizationStr, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing cpu utilization: %w", err)
	}

	return utilization, nil
}

func getMemoryUtilization() (float64, error) {
	cmd := exec.Command("./bin/memusage") // assuming the C program is compiled and present in the same directory.
	out, err := cmd.Output()
	if err != nil {
		return 0.0, fmt.Errorf("error executing memusage: %w", err)
	}

	// Assuming the output is in the format "Memory Utilization: xx.xx%"
	outputStr := strings.TrimSpace(string(out))
	splitStr := strings.Split(outputStr, ": ")
	if len(splitStr) < 2 {
		return 0.0, fmt.Errorf("unexpected output format: %s", outputStr)
	}

	utilizationStr := strings.TrimSuffix(splitStr[1], "%")
	utilization, err := strconv.ParseFloat(utilizationStr, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing memory utilization: %w", err)
	}

	return utilization, nil
}

func getDiskUtilization() (float64, error) {
	cmd := exec.Command("./bin/diskusage") // assuming the C program is compiled and present in the same directory.
	out, err := cmd.Output()
	if err != nil {
		return 0.0, fmt.Errorf("error executing diskusage: %w", err)
	}

	// Assuming the output is in the format "Disk Utilization: xx.xx%"
	outputStr := strings.TrimSpace(string(out))
	splitStr := strings.Split(outputStr, ": ")
	if len(splitStr) < 2 {
		return 0.0, fmt.Errorf("unexpected output format: %s", outputStr)
	}

	utilizationStr := strings.TrimSuffix(splitStr[1], "%")
	utilization, err := strconv.ParseFloat(utilizationStr, 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing disk utilization: %w", err)
	}

	return utilization, nil
}

func main() {
	if len(os.Args) < 2 {
		panic("please specify a port")
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
		fmt.Println("Hello, world!")
	})

	http.HandleFunc("/metrics/cpu", func(w http.ResponseWriter, r *http.Request) {
		utilization, err := getCPUUtilization()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting CPU utilization: %v", err), http.StatusInternalServerError)
			return
		}
		// Convert the float64 value to a string and then to []byte
		utilizationStr := fmt.Sprintf("%.2f%%", utilization)
		w.Write([]byte(utilizationStr))
		fmt.Printf("CPU Utilization: %s\n", utilizationStr)
	})

	http.HandleFunc("/metrics/memory", func(w http.ResponseWriter, r *http.Request) {
		utilization, err := getMemoryUtilization()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting memory utilization: %v", err), http.StatusInternalServerError)
			return
		}
		// Convert the float64 value to a string and then to []byte
		utilizationStr := fmt.Sprintf("%.2f%%", utilization)
		w.Write([]byte(utilizationStr))
		fmt.Printf("Memory Utilization: %s\n", utilizationStr)
	})

	http.HandleFunc("/metrics/disk", func(w http.ResponseWriter, r *http.Request) {
		utilization, err := getDiskUtilization()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting disk utilization: %v", err), http.StatusInternalServerError)
			return
		}
		// Convert the float64 value to a string and then to []byte
		utilizationStr := fmt.Sprintf("%.2f%%", utilization)
		w.Write([]byte(utilizationStr))
		fmt.Printf("Memory Utilization: %s\n", utilizationStr)
	})

	port := os.Args[1]

	fmt.Println("Server running on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
