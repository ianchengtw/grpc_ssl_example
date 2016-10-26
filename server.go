package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	pb "./defines"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	tlsCert = flag.String("tls-cert", "cert/cert.pem", "TLS server certificate.")
	tlsKey  = flag.String("tls-key", "cert/key.pem", "TLS server key.")
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *pb.Request) (*pb.Response, error) {

	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	token, ok := md["authorization"]
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return &pb.Response{Message: fmt.Sprintf("Hello %s (%s)", in.Name, token)}, nil
}

func NewServer() *server {
	return &server{}
}

func main() {
	flag.Parse()

	log.Println("Helloworld service starting...")

	cert, err := tls.LoadX509KeyPair(*tlsCert, *tlsKey)
	if err != nil {
		log.Fatal(err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	s := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterGreeterServer(s, NewServer())

	lis, _ := net.Listen("tcp", ":10000")
	s.Serve(lis)

}
