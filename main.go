package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reddit-clone/client"
	"reddit-clone/engine"
	"time"
)

type SimulationConfig struct {
	NumUsers    int `json:"numUsers"`
	NumSRs      int `json:"numSRs"`
	NumPosts    int `json:"numPosts"`
	NumComments int `json:"numComments"`
	NumVotes    int `json:"numVotes"`
	NumMessages int `json:"numMessages"`
}

type Config struct {
	Simulations []SimulationConfig `json:"simulations"`
}

func readConfig(filePath string) (*Config, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	return &config, nil
}

func main() {
	loggingFlag := flag.Bool("logging", true, "Enable or disable logging (true/false)")
	apiFlag := flag.Bool("api", false, "Start the REST API server")
	flag.Parse()

	if *apiFlag {
		startAPIServer(*loggingFlag)
	} else {
		runSimulation(*loggingFlag)
	}
}

func runSimulation(loggingEnabled bool) {
	configFile := "sim_config.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("Configuration file not found.")
		return
	}

	config, err := readConfig(configFile)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	fmt.Printf("Loaded %d simulations from config.\n", len(config.Simulations))
	for i, simConfig := range config.Simulations {
		client.LoggingEnabled = loggingEnabled
		if client.LoggingEnabled {
			fmt.Println("Logging is enabled.")
		} else {
			fmt.Println("Logging is disabled.")
		}

		fmt.Printf("\nRunning simulation #%d with parameters: %+v\n", i+1, simConfig)
		sim := client.NewSimulator()
		start := time.Now()
		sim.Run(simConfig.NumUsers, simConfig.NumSRs, simConfig.NumPosts, simConfig.NumComments, simConfig.NumVotes, simConfig.NumMessages)
		elapsed := time.Since(start)
		fmt.Printf("Simulation #%d completed in %s\n", i+1, elapsed)
		fmt.Printf("Users: %d\n", len(sim.Engine.Users))
		fmt.Printf("SubReddits: %d\n", len(sim.Engine.SubReddits))
		fmt.Printf("Messages: %d\n", len(sim.Engine.Messages))
		fmt.Println("-------------------------------")
	}
}

func startAPIServer(loggingEnabled bool) {
	redditEngine := engine.NewRedditEngine()
	api := NewAPI(redditEngine)

	if loggingEnabled {
		fmt.Println("Logging is enabled.")
	} else {
		fmt.Println("Logging is disabled.")
	}

	http.Handle("/api/", corsMiddleware(api))
	fmt.Println("REST API server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

// CORS Middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}