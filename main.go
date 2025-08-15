package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Sheikh-Fahad-Ahmed/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secretKey      string
	polkaKey       string
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

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		errorHandler(w, r, http.StatusForbidden, "", nil)
		return
	}
	cfg.db.DeleteAllUsers(r.Context())
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "ok")
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Error connection got database")
		os.Exit(1)
	}
	dbQueries := database.New(dbConn)

	apicfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
		secretKey:      os.Getenv("SECRET_KEY"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	const port = "8080"
	mux := http.NewServeMux()

	mux.Handle("/app/", apicfg.middlewareMetricInc(
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	))

	mux.HandleFunc("GET /admin/metrics", apicfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apicfg.resetHandler)

	mux.HandleFunc("GET /api/healthz", myHandler)
	mux.HandleFunc("GET /api/chirps", apicfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.getChirp)
	mux.HandleFunc("POST /api/login", apicfg.authenticateUser)
	mux.HandleFunc("POST /api/chirps", apicfg.chirpHandler)
	mux.HandleFunc("POST /api/users", apicfg.userHandler)
	mux.HandleFunc("POST /api/refresh", apicfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", apicfg.revokeHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apicfg.polkaWebhookHandler)

	mux.HandleFunc("PUT /api/users", apicfg.userPUTHandler)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apicfg.deleteChirpHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Started a server at port:", port)
	log.Fatal(server.ListenAndServe())
}
