package handler

import (
	"log"
	"errors"
	"golang.org/x/net/context"
	"github.com/henter/go-zipkin-demo/pb"
	"strings"
)

type Demo struct {}

func (d *Demo) Hello(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
	log.Print("Received Demo.Hello request")
	rsp.Code = 0
	rsp.Msg = "hello world"

	if req.Q == "" || req.N <= 0 || req.N > 100 {
		return errors.New("parameter error")
	}

	// just repeat N times
	rsp.Msg = strings.Repeat(req.Q, int(req.N))

	return nil
}
