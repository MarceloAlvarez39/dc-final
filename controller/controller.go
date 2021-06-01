package controller

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.nanomsg.org/mangos"

	// register transports

	"go.nanomsg.org/mangos/protocol/surveyor"
	_ "go.nanomsg.org/mangos/transport/all"
)

var controllerAddress = "tcp://localhost:40899"
var sock mangos.Socket

// ----------------------- STRUCTURES -----------------------------
// ----------------------------------------------------------------

type Workload struct {
	WorkloadID     string
	Filter         string
	WorkloadName   string
	Status         string
	RunningJobs    int
	FilteredImages []string
}

type Worker struct {
	Name     string
	Tags     string
	Status   string
	Usage    int
	URL      string
	Active   bool
	Port     int
	JobsDone int
}

// ------------------------ CODE ---------------------
// ---------------------------------------------------

var Workloads = make(map[string]Workload)
var Workers = make(map[string]Worker)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Start() {
	var sock mangos.Socket
	var err error
	if sock, err = surveyor.NewSocket(); err != nil { // We used the survey method
		die("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}

	for {
		// Could also use sock.RecvMsg to get header
		d := date()
		log.Printf("Controller: Publishing Date %s\n", d)
		if err = sock.Send([]byte(d)); err != nil {
			die("Failed publishing: %s", err.Error())
		}
		time.Sleep(time.Second * 3)
	}
}
