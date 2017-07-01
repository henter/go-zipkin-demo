package main

import (
	"os"
	"log"

	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	"go-zipkin-demo/pb"
	"go-zipkin-demo/handler"
	"go-zipkin-demo/trace"
)


func main() {
	service_name := "go.zipkin.demo"

	zipkin_addr := "http://localhost:9411/api/v1/spans"
	consul_addr := "127.0.0.1:8500"
	hostname, _ := os.Hostname()
	InitTracer(zipkin_addr, hostname, service_name)

	reg := consul.NewRegistry(
		registry.Addrs(consul_addr),
	)
	service := grpc.NewService(
		micro.Name(service_name),
		micro.Version("v0.1"),
		micro.Registry(reg),
		micro.WrapHandler(trace.ServerWrapper),
	)

	service.Init()

	pb.RegisterDemoHandler(service.Server(), new(handler.Demo))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatalf("service error: %s", err.Error())
	}
}

func InitTracer(zipkinURL string, hostPort string, serviceName string) {
	collector, err := zipkin.NewHTTPCollector(zipkinURL)
	if err != nil {
		log.Fatalf("unable to create Zipkin HTTP collector: %v", err)
		return
	}
	tracer, err := zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, hostPort, serviceName),
	)
	if err != nil {
		log.Fatalf("unable to create Zipkin tracer: %v", err)
		return
	}
	opentracing.InitGlobalTracer(tracer)
	return
}


