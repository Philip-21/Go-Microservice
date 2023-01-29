package main

import (
	"context"
	"log"
	"log-service/database"
	"time"
)

type RPCServer struct {
}

// Values to be inserted in Db
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), database.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error writing t mongo", err)
		return err
	}
	//a pointer to the string to send a message back to the
	//person who called it
	*resp = "Processed payload Via RPC:" + payload.Name
	return nil
}
