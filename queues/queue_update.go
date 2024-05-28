// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package queues

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/wade-rees-me/striker-go/constants"
)

func Update(queue string, timeout int64, receiptHandle string) {
	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	urlResult, err := GetQueueURL(sess, &queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	changeParams := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(*urlResult.QueueUrl),
		ReceiptHandle:     &receiptHandle,
		VisibilityTimeout: aws.Int64(constants.QueueTimeout),
	}

	svc := sqs.New(sess)
	_, err = svc.ChangeMessageVisibility(changeParams)
	if err != nil {
		log.Fatalf("Unable to change visibility timeout for message %v, %v.", receiptHandle, err)
	}
	//fmt.Println("Visibility timeout for message updated successfully.")
}

/*
func main() {
    changeParams := &sqs.ChangeMessageVisibilityInput{
        QueueUrl:          aws.String(queueURL),
        ReceiptHandle:     result.Messages[0].ReceiptHandle,
        VisibilityTimeout: aws.Int64(60), // New visibility timeout in seconds
    }

    _, err = svc.ChangeMessageVisibility(changeParams)
    if err != nil {
        log.Fatalf("Unable to change visibility timeout for message %v, %v.", *result.Messages[0].MessageId, err)
    }

    fmt.Println("Visibility timeout for message updated successfully.")

    // Example: Send a message with a request timeout
    sendParams := &sqs.SendMessageInput{
        QueueUrl:    aws.String(queueURL),
        MessageBody: aws.String("This is a test message"),
    }

    // Setting a timeout for the send message request
    svc.SendMessageWithContext(aws.ContextWithTimeout(sess, 5*time.Second), sendParams)

    if err != nil {
        log.Fatalf("Unable to send message to queue %q, %v.", queueURL, err)
    }

    fmt.Println("Message sent successfully.")
}

*/
