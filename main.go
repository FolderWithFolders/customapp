package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var rtp float64
var maxMultiplier = 10000.0

func main() {
	flag.Float64Var(&rtp, "rtp", 1.0, "RTP value in (0, 1.0]")
	flag.Parse()

	if rtp <= 0 || rtp > 1.0 {
		log.Fatal("rtp must be in (0, 1.0]")
	}

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/get", getHandler)

	log.Println("Server starting on :64333 with RTP:", rtp)
	log.Fatal(http.ListenAndServe(":64333", nil))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	multiplier := generateMultiplier(rtp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"result": multiplier})
}

func generateMultiplier(rtp float64) float64 {
	if rand.Float64() < 1-rtp {
		return 1.0
	}

	// Генерация unbounded Pareto с capping
	v := rand.Float64()
	denom := 1.0 - v
	if denom <= 0 || denom < 1.0/maxMultiplier { // Избежать inf и cap
		return maxMultiplier
	}
	m := 1.0 / denom
	if m > maxMultiplier {
		return maxMultiplier
	}
	return m
}
