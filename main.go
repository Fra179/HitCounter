package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ctx = context.Background()
	rdc *redis.Client

	// this regex is not used anymore
	// regex taken from https://regex101.com/library/SEg6KL
	// it has some imperfections, but it's good enough
	// validateDomain = regexp.MustCompile(`(?im)^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)

	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisHost     = os.Getenv("REDIS_HOST")
	redisPort     = os.Getenv("REDIS_PORT")
	appPort       = os.Getenv("APP_PORT")
	appHost       = os.Getenv("APP_HOST")
)

const TIMEOUT int64 = 60

/*
	This code was written to generate standardized json responses for
	the APIs. It's useless right now because of the simple nature of the
	program. I've commented it out because I might need it in the future
*/
//type ResponseObject struct {
//	StatusCode  int
//	ContentType string
//	Message     string
//}
//
//func (r ResponseObject) GenerateResponse(w http.ResponseWriter) {
//	w.WriteHeader(r.StatusCode)
//	w.Write([]byte(r.Body))
//}

func GenerateCookie(w http.ResponseWriter, url string) {
	http.SetCookie(w, &http.Cookie{
		Name:  url,
		Value: strconv.Itoa(int(time.Now().Unix())),
	})
}

func CheckCookie(c *http.Cookie, _ error) bool {
	if c == nil {
		return true
	}

	val, err := strconv.Atoi(c.Value)

	return err != nil || time.Now().Unix()-int64(val) > TIMEOUT
}

func Increment(url string) int {
	if res, err := rdc.Incr(ctx, url).Result(); err == nil {
		return int(res)
	}

	if _, err := rdc.Set(ctx, url, 1, 0).Result(); err == nil {
		return 1
	}

	return 0
}

func Get(url string) int {
	if res, err := rdc.Get(ctx, url).Result(); err == nil {
		if val, err := strconv.Atoi(res); err == nil {
			return val
		}
	}

	return 0
}

func HitCounter(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(r.URL.Query().Get("url"))

	if url == "" {
		fmt.Fprintf(w, "0")
		return
	}

	var count int

	if CheckCookie(r.Cookie(url)) {
		count = Increment(url)
		GenerateCookie(w, url)
	} else {
		count = Get(url)
	}

	fmt.Fprint(w, count)
}

func HitCounterSilent(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(r.URL.Query().Get("url"))

	if url == "" {
		fmt.Fprintf(w, "0")
		return
	}

	fmt.Fprint(w, Get(url))
}

func ParseEnv() {
	toCheck := map[*string]string{
		&redisHost:     "localhost",
		&redisPort:     "6379",
		&appHost:       "0.0.0.0",
		&appPort:       "8080",
		&redisPassword: "",
	}

	for k, v := range toCheck {
		if *k == "" {
			*k = v
		}
	}
}

func CorsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rdc.Incr(ctx, "totalHits")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")
		h.ServeHTTP(w, r)
	})
}

func main() {
	ParseEnv()
	rdc = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		DB:       0, // use default DB
		Password: redisPassword,
	})

	if err := rdc.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/count", HitCounter)
	mux.HandleFunc("/get", HitCounterSilent)
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
	})

	fmt.Printf("Started server on %s:%s\n", appHost, appPort)
	log.Fatal(http.ListenAndServe(appHost+":"+appPort, CorsMiddleware(mux)))
}
