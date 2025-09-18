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
var rng *rand.Rand

func main() {
	flag.Float64Var(&rtp, "rtp", 1.0, "RTP value in (0, 1.0]")
	flag.Parse()

	if rtp <= 0 || rtp > 1.0 {
		log.Fatal("rtp must be in (0, 1.0]")
	}

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	if rng.Float64() < rtp {
		return 1.0 + rng.Float64()*(maxMultiplier-1.0)
	} else {
		return rng.Float64() * 1.0
	}
}
