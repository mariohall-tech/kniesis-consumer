package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/harlow/kinesis-connectors"
	"github.com/joho/godotenv"
	"github.com/sendgridlabs/go-kinesis"
)

func main() {
	godotenv.Load()

	// Initialize Kinesis client
	auth := kinesis.NewAuth()
	ksis := kinesis.New(&auth, kinesis.Region{})

	// Create stream
	connector.CreateStream(ksis, "userStream", 2)

	// read file
	// https://s3.amazonaws.com/kinesis.test/users.txt
	file, _ := os.Open("tmp/users.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	args := kinesis.NewArgs()
	args.Add("StreamName", "userStream")
	ctr := 0

	for scanner.Scan() {
		l := scanner.Text()
		ctr = ctr + 1
		key := fmt.Sprintf("partitionKey-%d", ctr)

		args := kinesis.NewArgs()
		args.Add("StreamName", "userStream")
		args.AddRecord([]byte(l), key)
		ksis.PutRecords(args)
		fmt.Print(".")
	}

	fmt.Println(".")
	fmt.Println("Finished populating userStream")
}
