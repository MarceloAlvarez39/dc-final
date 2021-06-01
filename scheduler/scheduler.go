package scheduler

import (
	"context"
	"log"
	"time"

	pb "github.com/marceloalvarez39/dc-final/proto"
	"google.golang.org/grpc"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

//var schedulerAddress = "50051"

// ----------------------- STRUCTURES -----------------------------
// ----------------------------------------------------------------

type Job struct {
	Address string
	RPCName string
	Filter  string
	ImageID string
	Worker  string //due to change to another type
}

// ----------------------- CODE -----------------------------
// ----------------------------------------------------------------

func schedule(job Job) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFiltersClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	r, err := c.GrayScale(ctx, &pb.FilterRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Scheduler: RPC respose from %s : %s", job.Address, r.GetMessage())

	m, err := c.Blur(ctx, &pb.FilterRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Scheduler: RPC respose from %s : %s", job.Address, m.GetMessage())
}

func Start(jobs chan Job) error {
	for {
		job := <-jobs
		schedule(job)
	}

}
