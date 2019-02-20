package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/abicky/ecswrap/internal/log"
	"github.com/jessevdk/go-flags"
)

type options struct {
	Timeout               uint     `long:"stop-wait-timeout" env:"ECSWRAP_STOP_WAIT_TIMEOUT" default:"10" description:"Maximum time duration in seconds to wait from when the process receives SIGTERM before sending SIGTERM to the child. This value should be less than ECS_CONTAINER_STOP_TIMEOUT."`
	LinkedContainers      []string `long:"linked-container" env:"ECSWRAP_LINKED_CONTAINERS" env-delim:"," description:"container names linked with the container where this program is running."`
	SignalForwardingDelay uint     `long:"signal-forwarding-delay" env:"ECSWRAP_SIGNAL_FORWARDING_DELAY" default:"0" description:"Delay seconds until forwarding a signal, which is SIGTERM, SIGQUIT or SIGINT,to child processes."`
	Verbosity             []bool   `short:"v" long:"verbose" description:"Verbosity"`
}

type childProcess struct {
	cmd        *exec.Cmd
	exitStatus int
}

type task struct {
	Containers []*container
}

type container struct {
	Name          string
	DesiredStatus string
	KnownStatus   string
}

const defaultErrorExitStatus = 1

// URL of the task metadata version 3 endpoint
// cf. https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v3.html
var ECSMetaDataURL = os.ExpandEnv("${ECS_CONTAINER_METADATA_URI}/task")

var forwardableSignals = []os.Signal{
	syscall.SIGCONT,
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
	syscall.SIGTSTP,
	syscall.SIGUSR1,
	syscall.SIGUSR2,
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] -- COMMAND [ARGS]"

	args, err := parser.Parse()
	if err != nil || len(args) == 0 {
		if !flags.WroteHelp(err) {
			parser.WriteHelp(os.Stderr)
		}
		os.Exit(defaultErrorExitStatus)
	}

	log.SetOutput(os.Stdout)
	log.SetPrefix(parser.Name + " ")
	log.SetVerbosity(len(opts.Verbosity))

	log.Debugf("opts: %+v, args: %v, PID: %v\n", opts, args, os.Getpid())
	child, err := start(args, opts.LinkedContainers, opts.Timeout, opts.SignalForwardingDelay)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(defaultErrorExitStatus)
	}

	exitCode, err := child.wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

func start(args []string, linkedContainers []string, timeout uint, signalForwardingDelay uint) (*childProcess, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Failed to start command %v: %v\n", args, err)
	}
	log.Debugf("Succeeded to start command %v with PID %d\n", args, cmd.Process.Pid)

	go handleSignals(cmd.Process, linkedContainers, timeout, signalForwardingDelay)

	return &childProcess{cmd: cmd}, nil
}

func (cp *childProcess) wait() (int, error) {
	err := cp.cmd.Wait()
	if err == nil {
		log.Debugf("child process exited normally")
		return 0, nil
	}

	status := defaultErrorExitStatus
	if exitError, ok := err.(*exec.ExitError); ok {
		if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if ws.Exited() {
				status = ws.ExitStatus()
				log.Debugf("child process exited with status %d\n", status)
			} else if ws.Signaled() {
				log.Debugf("child process was interrupted by signal %d\n", ws.Signal())
				status = 128 + int(ws.Signal())
			} else {
				log.Warnf("child process exited with unknown wait status %v", ws)
			}
		}
	}

	return status, err
}

func handleSignals(process *os.Process, linkedContainers []string, timeout uint, signalForwardingDelay uint) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, forwardableSignals...)

	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
			waitContainers(linkedContainers, timeout)
			if signalForwardingDelay > 0 {
				time.Sleep(time.Duration(signalForwardingDelay) * time.Second)
			}
		}
		log.Debugf("Send signal %d to child process\n", sig)
		process.Signal(sig)
	}
}

func waitContainers(linkedContainers []string, timeout uint) {
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			log.Warnf("LinkedContainers didn't stop within %d seconds\n", timeout)
			return
		case <-ticker.C:
			task, err := getECSTask()
			if err != nil {
				log.Warnln("Failed to get ECS task information:", err)
				continue
			}

			allStopped := true
			for _, c := range task.Containers {
				log.Debugf("container: %+v\n", c)
				if contains(linkedContainers, c.Name) && c.KnownStatus != "STOPPED" {
					allStopped = false
					break
				}
			}
			if allStopped {
				log.Debugln("All linked containers have stopped")
				return
			}
		}
	}
}

func getECSTask() (*task, error) {
	resp, err := http.Get(ECSMetaDataURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var task task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func contains(strArr []string, targetStr string) bool {
	for _, str := range strArr {
		if str == targetStr {
			return true
		}
	}
	return false
}
