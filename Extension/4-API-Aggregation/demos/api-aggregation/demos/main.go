package main

import (
    "context"
    "crypto/tls"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/go-redis/redis/v8"
)

var (
    ctx = context.Background()
    rdb *redis.Client
)

type Vote struct {
    Vote    string `json:"vote"`
    VoterID string `json:"voter_id"`
}

func initRedis() {
    redisAddr := os.Getenv("REDIS_URL")
    if redisAddr == "" {
        redisAddr = "redis://redis:6379"
    }

    opt, err := redis.ParseURL(redisAddr)
    if err != nil {
        log.Fatalf("Error parsing Redis URL: %v", err)
    }

    rdb = redis.NewClient(opt)
}

func handleGetVotes(w http.ResponseWriter, r *http.Request) {
    votes, err := rdb.LRange(ctx, "votes", 0, -1).Result()
    if err != nil {
        http.Error(w, "error retrieving votes", http.StatusInternalServerError)
        return
    }

    var result []Vote
    for _, voteData := range votes {
        var vote Vote
        if err := json.Unmarshal([]byte(voteData), &vote); err != nil {
            http.Error(w, "error unmarshalling vote data", http.StatusInternalServerError)
            return
        }
        result = append(result, vote)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func main() {
    initRedis()

    mux := http.NewServeMux()
    mux.HandleFunc("/apis/mycompany.com/v1/votes", handleGetVotes)

    server := &http.Server{
        Addr:         ":8443",
        Handler:      mux,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }

    log.Println("Starting server on :8443")
    log.Fatal(server.ListenAndServeTLS("/tls/tls.crt", "/tls/tls.key"))
}
