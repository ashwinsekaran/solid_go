package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Notification struct {
	Details  NotificationDetails
	Response chan Result
	Ctx      context.Context
}

type NotificationDetails struct {
	Message   string
	Recipient string
}

type Result struct {
	Success string
	Err     error
}

func main() {
	wg := sync.WaitGroup{}

	jobCh := make(chan Notification, 5)

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				processNotification(job)
			}
		}()
	}

	r := httprouter.New()

	r.Handle("POST", "/notify", notification(jobCh))

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, os.Interrupt)

	<-shutdown
	close(jobCh)
	wg.Wait()

}

func notification(jobCh chan<- Notification) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var notify Notification
		var data NotificationDetails

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		notify.Response = make(chan Result, 1)
		notify.Ctx = ctx

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		notify.Details = data

		jobCh <- notify

		result := <-notify.Response
		if result.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result.Success))

	}

}

func processNotification(n Notification) {
	select {
	case <-time.After(2 * time.Second):
		res := fmt.Sprintf("notification message: %s for receipeint: %s", n.Details.Message, n.Details.Recipient)
		n.Response <- Result{Success: res}
	case <-n.Ctx.Done():
		n.Response <- Result{Err: n.Ctx.Err()}
	}
}
