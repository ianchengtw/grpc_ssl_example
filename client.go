package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	pb "./defines"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type IdentifiedToken struct {
	Token string
}

func (it *IdentifiedToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": it.Token,
	}, nil
}

func (it *IdentifiedToken) RequireTransportSecurity() bool {
	return true
}

func main() {
	var (
		caCert = flag.String("ca", "cert/cert.pem", "Trusted CA certificate.")
		server = flag.String("server", ":10000", "Server address.")
		name   = flag.String("name", "ian", "Username to use.")
		token  = "testing_token"
	)
	flag.Parse()

	rawCACert, err := ioutil.ReadFile(*caCert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(rawCACert)

	creds := credentials.NewTLS(&tls.Config{
		RootCAs: caCertPool,
		// ServerName: "secure-server.example.com",
		InsecureSkipVerify: true,
	})

	conn, err := grpc.Dial(*server,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&IdentifiedToken{Token: token}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	response, err := c.SayHello(context.Background(), &pb.Request{Name: *name})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response: ", response.Message)

}
