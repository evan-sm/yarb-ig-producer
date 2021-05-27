package main

import (
	//	"time"

	"context"
	"encoding/json"
	"fmt"
	//	"google.golang.org/grpc/codes"
	//	"google.golang.org/grpc/status"
	"cloud.google.com/go/pubsub"
	log "github.com/sirupsen/logrus"
	yarb "github.com/wmw9/yarb-struct"
)

func SendToPubSub(projectID, topicID string, p yarb.Payload) error {
	msg, _ := json.Marshal(p) // Convert payload to json bytes

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}

	log.Printf("Published a message to %v, %v; msg ID: %v\n", projectID, topicID, id)
	return nil
}
