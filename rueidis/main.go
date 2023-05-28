package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/redis/rueidis"
)

func main() {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := client.Do(r.Context(), client.B().Ping().Build()).Error()
		if err != nil {
			log.Printf("failed to ping: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("connected\n"))
	})
	http.HandleFunc("/value", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHandler(client, w, r)
		case http.MethodPost:
			postHandler(client, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func getHandler(client rueidis.Client, w http.ResponseWriter, r *http.Request) {
	defer func(now time.Time) { log.Printf("getHandler took %s\n", time.Since(now)) }(time.Now())

	get, err := get(r.Context(), client, "key")
	if err != nil {
		log.Printf("failed to get value: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(get + "\n"))
}

func get(ctx context.Context, client rueidis.Client, key string) (string, error) {
	cmd := client.B().Get().Key(key).Cache()
	result := client.DoCache(ctx, cmd, time.Minute)
	err := result.Error()
	if err != nil {
		return "", err
	}
	return result.ToString()
}

func postHandler(client rueidis.Client, w http.ResponseWriter, r *http.Request) {
	defer func(now time.Time) { log.Printf("postHandler took %s\n", time.Since(now)) }(time.Now())

	v := r.URL.Query().Get("kv")
	if v == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("kv is empty\n"))
		return
	}
	log.Printf("set value: %s\n", v)
	err := set(r.Context(), client, "key", v)
	if err != nil {
		log.Printf("failed to set value: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("success\n"))
}

func set(ctx context.Context, client rueidis.Client, key, value string) error {
	cmd := client.B().Set().
		Key(key).Value(value).
		// Nx().
		Build()
	err := client.Do(ctx, cmd).Error()
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			return err
		}
		log.Printf("key(%s) already exists", key)
	}
	return nil
}
