package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, r)
	})
}

func myHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, req *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<html>
  							<body>
    							<h1>Welcome, Chirpy Admin</h1>
    							<p>Chirpy has been visited %d times!</p>
  							</body>
						</html>
						`, hits)
	w.Write([]byte(html))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "ok")
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type respParams struct {
		Body string `json:"body"`
	}

	type errorRespParams struct {
		Error string `json:"error"`
	}

	type ValidRespParams struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := respParams{}
	err := decoder.Decode(&params)
	if err != nil {
		errResponse := errorRespParams{
			Error: "Something went wrong",
		}

		data, err := json.Marshal(errResponse)
		if err != nil {
			log.Printf("Error marshaling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	if len(params.Body) > 140 {
		errResponse := errorRespParams{
			Error: "Chirp is too long",
		}

		data, err := json.Marshal(errResponse)
		if err != nil {
			log.Printf("Error marshaling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
	} else {
		validResponse := ValidRespParams{
			Valid: true,
		}

		data, err := json.Marshal(validResponse)
		if err != nil {
			log.Printf("Error marshaling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data)
	}
}

func main() {

	apicfg := apiConfig{}

	const port = "8080"
	mux := http.NewServeMux()

	mux.Handle("/app/", apicfg.middlewareMetricInc(
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	))

	
	
	mux.HandleFunc("GET /admin/metrics", apicfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.resetHandler)

	mux.HandleFunc("GET /api/healthz", myHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Started a server at port:", port)
	log.Fatal(server.ListenAndServe())
}
