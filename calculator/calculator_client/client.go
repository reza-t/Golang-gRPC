package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/reza-t/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// doClientStreaming(c)

	// doBiDiStreaming(c)

	doErrorUnary(c)

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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("staring to do a doClientStreaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v\n", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream: %v\n", err)
		return
	}

	waitC := make(chan struct{})

	go func() {
		for {
			err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: 1,
			})

			if err == io.EOF {
				log.Fatalf("End of the file")

				break
			}

			if err != nil {
				break
			}
		}
		close(waitC)
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				break
			}
			fmt.Println(res.GetNumber())
		}
		close(waitC)
	}()

	<-waitC
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("staring to do a doErrorUnary RPC...")

	//correct call
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: 10})

	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC (user error)
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent sent a negative number")
			}
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("The response for correct request is: %v", res)

	//wrong call
	noRes, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: -10})

	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC (user error)
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent sent a negative number")
			}
			return
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("The response for correct request is: %v", noRes)
}
