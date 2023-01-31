package main

import (
	"context"
	"fmt"
	"log"
	"log-service/database"
	"net"
	"net/rpc"
	"time"
)

type RPCServer struct {
}

// Values to be inserted in Db
type RPCPayload struct {
	Name string
	Data string
}

// writes a response back to the broker after its saved in the db
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

func (app *Config) rpcListen() error {
	log.Printf("Starting RPC server on Port %s", rpcPort)
	//Listen announces on the local network address.
	//The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Println(err)
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
			//start over the connection
		}
		go rpc.ServeConn(rpcConn)
	}
}
