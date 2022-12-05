package server

import (
	h "b0b-common/internal/server/proto"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"testing"
)

func TestHelloServer(t *testing.T) {
	//stream
	conn, err := grpc.Dial("127.0.0.1:8889", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := h.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &h.HelloRequest{Name: "bobby"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}
