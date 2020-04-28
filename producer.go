package main

import (
	"context"
	"fmt"
	"os"

	_ "math/bits"
	"time"
	"pack.ag/amqp"

//import "pack.ag/amqp"
uuid "github.com/satori/go.uuid"
)

func failOnError(err error, message string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("%s: %s", message, err))
		os.Exit(1)
	}
}

func main (){
// Set custom env variable
	var username = os.Getenv("username")
	var password = os.Getenv("password")
	var host = os.Getenv("host")
	var port = os.Getenv("port")

	// Create session
	client, err := amqp.Dial("amqps://"+host+":"+port,
		amqp.ConnSASLPlain(username, password),
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("Dialing AMQP server:", err))
	}
	defer client.Close()

	// Open a session
	session, err := client.NewSession()
	ctx := context.Background()

	queue := "/test-queue"

	sender, err := session.NewSender(
		amqp.LinkTargetAddress(queue),
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot create session, error:", err))
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
//	Sending messages to queue:

	totalSeconds := 0
	maxSeconds := 30
	for totalSeconds < maxSeconds {
		// Generate random UUID for message body. Convert to []byte
		randomUUID := uuid.NewV4()
		body := []byte(randomUUID.String())

		// Publish a message
		output := fmt.Sprintf("Published message to test-queue: %s", body)
		fmt.Println(output)

		err = sender.Send(ctx, amqp.NewMessage([]byte(body)))
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to publish message, error:", err))
		}
		time.Sleep(5 * time.Second)
		totalSeconds += 5
	}

	// Close connection
	sender.Close(ctx)
	cancel()
}
