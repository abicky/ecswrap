package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/abicky/ecswrap/internal/log"
)

type metadataServer struct {
	port uint
	task *task
}

func startMetadataServer(port uint, task *task) {
	ms := &metadataServer{port: port, task: task}
	http.HandleFunc("/task", ms.handler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func (ms *metadataServer) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(ms.task)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(500)
	}
	fmt.Fprintf(w, string(jsonBytes))
}

func TestStart(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	linkedContainer := &container{Name: "linked", DesiredStatus: "RUNNING", KnownStatus: "RUNNING"}
	notLinkedContainer := &container{Name: "not-linked", DesiredStatus: "RUNNING", KnownStatus: "RUNNING"}
	task := &task{Containers: []*container{linkedContainer, notLinkedContainer}}

	go startMetadataServer(8080, task)

	timeout := uint(3)

	t.Run("it exits with status 0", func(t *testing.T) {
		cp, err := start([]string{"sleep", "0"}, []string{"linked"}, timeout, 0)
		if err != nil {
			t.Fatal(err)
		}

		exitStatus := make(chan int)
		go func() {
			want := 0
			select {
			case got := <-exitStatus:
				if got != want {
					t.Errorf("exitStatus: got: %v, want: %v", got, want)
				}
				return
			case <-time.After(time.Duration(timeout) * time.Second):
				t.Errorf("timeout")
			}
		}()

		status, err := cp.wait()
		if err != nil {
			t.Errorf("err is expected to be nil")
		}
		exitStatus <- status
	})

	t.Run("it sends SIGUSR1 immediately", func(t *testing.T) {
		cp, err := start([]string{"sleep", "10"}, []string{"linked"}, timeout, 0)
		if err != nil {
			t.Fatal(err)
		}

		exitStatus := make(chan int)
		go func() {
			process, _ := os.FindProcess(os.Getpid())
			process.Signal(syscall.SIGUSR1)

			want := 128 + int(syscall.SIGUSR1)
			select {
			case got := <-exitStatus:
				if got != want {
					t.Errorf("exitStatus: got: %v, want: %v", got, want)
				}
				return
			case <-time.After(time.Duration(timeout) * time.Second):
				t.Errorf("timeout")
			}
		}()

		status, err := cp.wait()
		if err == nil {
			t.Errorf("err is expected to be not nil")
		}
		exitStatus <- status
	})

	t.Run("it sends SIGTERM after timeout", func(t *testing.T) {
		cp, err := start([]string{"sleep", "10"}, []string{"linked"}, timeout, 0)
		if err != nil {
			t.Fatal(err)
		}

		exitStatus := make(chan int)
		go func() {
			process, _ := os.FindProcess(os.Getpid())
			process.Signal(syscall.SIGTERM)
			notLinkedContainer.KnownStatus = "STOPPED"

			want := 128 + int(syscall.SIGTERM)
			select {
			case <-exitStatus:
				t.Errorf("not timed out")
			case <-time.After(time.Duration(timeout-1) * time.Second):
			}

			got := <-exitStatus
			if got != want {
				t.Errorf("exitStatus: got: %v, want: %v", got, want)
			}
		}()

		status, err := cp.wait()
		if err == nil {
			t.Errorf("err is expected to be not nil")
		}
		exitStatus <- status
	})

	t.Run("it sends SIGTERM with signalForwardingDelay", func(t *testing.T) {
		signalForwardingDelay := uint(2)
		cp, err := start([]string{"sleep", "10"}, []string{"linked"}, timeout, signalForwardingDelay)
		if err != nil {
			t.Fatal(err)
		}

		exitStatus := make(chan int)
		go func() {
			process, _ := os.FindProcess(os.Getpid())
			process.Signal(syscall.SIGTERM)
			linkedContainer.KnownStatus = "STOPPED"

			want := 128 + int(syscall.SIGTERM)
			select {
			case got := <-exitStatus:
				if got != want {
					t.Errorf("exitStatus: got: %v, want: %v", got, want)
				}
				return
			case <-time.After(time.Duration(timeout+signalForwardingDelay+1) * time.Second):
				t.Errorf("timeout")
			}
			close(exitStatus)
		}()

		status, err := cp.wait()
		if err == nil {
			t.Errorf("err is expected to be not nil")
		}
		exitStatus <- status
	})

	t.Run("it sends SIGTERM without timeout", func(t *testing.T) {
		cp, err := start([]string{"sleep", "10"}, []string{"linked"}, timeout, 0)
		if err != nil {
			t.Fatal(err)
		}

		exitStatus := make(chan int)
		go func() {
			process, _ := os.FindProcess(os.Getpid())
			process.Signal(syscall.SIGTERM)
			linkedContainer.KnownStatus = "STOPPED"

			want := 128 + int(syscall.SIGTERM)
			select {
			case got := <-exitStatus:
				if got != want {
					t.Errorf("exitStatus: got: %v, want: %v", got, want)
				}
				return
			case <-time.After(time.Duration(timeout-1) * time.Second):
				t.Errorf("timeout")
			}
		}()

		status, err := cp.wait()
		if err == nil {
			t.Errorf("err is expected to be not nil")
		}
		exitStatus <- status
	})

}
