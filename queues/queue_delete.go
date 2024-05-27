// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package queues

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// DeleteMessage deletes a message from an Amazon SQS queue
// Inputs:
//
//	sess is the current session, which provides configuration for the SDK's service clients
//	queueURL is the URL of the queue
//	messageID is the ID of the message
//
// Output:
//
//	If success, nil
//	Otherwise, an error from the call to DeleteMessage
func DeleteMessage(sess *session.Session, queueURL *string, messageHandle *string) error {
	svc := sqs.New(sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	})
	if err != nil {
		return err
	}

	return nil
}

func Delete(queue, receiptHandle string) error {
	// Create a session that gets credential values from ~/.aws/credentials
	// and the default region from ~/.aws/config
	// snippet-start:[sqs.go.delete_message.sess]
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Get URL of queue
	result, err := GetQueueURL(sess, &queue)
	if err != nil {
		//fmt.Println("Got an error getting the queue URL:")
		//fmt.Println(err)
		return err
	}

	queueURL := result.QueueUrl
	err = DeleteMessage(sess, queueURL, &receiptHandle)
	if err != nil {
		//fmt.Println("Got an error deleting the message:")
		//fmt.Println(err)
		return err
	}

	//fmt.Println("Deleted message from queue with URL " + *queueURL)
	//fmt.Println("Deleted message from queue with Handle " + receiptHandle)
	return nil
}
