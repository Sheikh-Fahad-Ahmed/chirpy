package main

import (
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
