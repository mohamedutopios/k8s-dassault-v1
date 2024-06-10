package main

import (
    "context"
    "crypto/tls"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
)

var (
    ctx = context.Background()
    redisClient *redis.Client
)

type Vote struct {
    Vote    string `json:"vote"`
    VoterID string `json:"voter_id"`
}

func main() {
    // Initialize Redis client
    redisAddr := os.Getenv("REDIS_URL")
    if redisAddr == "" {
        redisAddr = "redis:6379"
    }

    redisClient = redis.NewClient(&redis.Options{
        Addr: redisAddr,
        TLSConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
    })

    mux := http.NewServeMux()
    mux.HandleFunc("/apis/mycompany.com/v1/votes", votesHandler)

    server := &http.Server{
        Addr:         ":8443",
        Handler:      mux,
        TLSConfig:    &tls.Config{MinVersion: tls.VersionTLS12},
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    log.Println("Starting server on :8443")
    log.Fatal(server.ListenAndServeTLS("/tls/tls.crt", "/tls/tls.key"))
}

func votesHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        handleGetVotes(w, r)
    case http.MethodPost:
        handlePostVote(w, r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleGetVotes(w http.ResponseWriter, r *http.Request) {
    // Retrieve votes from Redis
    votes, err := redisClient.LRange(ctx, "votes", 0, -1).Result()
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

func handlePostVote(w http.ResponseWriter, r *http.Request) {
    var vote Vote
    if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
        http.Error(w, "invalid request payload", http.StatusBadRequest)
        return
    }

    voteData, err := json.Marshal(vote)
    if err != nil {
        http.Error(w, "error marshalling vote data", http.StatusInternalServerError)
        return
    }

    if err := redisClient.RPush(ctx, "votes", voteData).Err(); err != nil {
        http.Error(w, "error saving vote", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "vote recorded"}`))
}
