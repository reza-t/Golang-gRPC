package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/reza-t/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client starting")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v\n", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("staring to do a Unary RPC...")
	num := make([]float64, 2)
	num[0] = 1.2
	num[1] = 2.4

	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			Numbers: num,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error, while calling Calculator RPC: %v", err)
	}
	log.Printf("Response from Calculator: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("staring to do a PrimeServerStreaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 9999,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error, while calling PrimeServerStreaming RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happen: %v", err)
		}

		fmt.Println(res.GetPrimeFactor())
	}
}
