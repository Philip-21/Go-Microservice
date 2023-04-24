The Logger Service 
- The Logger service receives and accepts the data from the authentication, listener and mailer service through the broker service and stores the data in a  mongoDb database
- Its responsible for handling the grpc and rpc actions 
- It sends a response back to the broker service after performing db actions. 
- It has an Api , Grpc , rpc and RabbitMQ response

- Uses  [chi router](https://github.com/go-chi/chi/v5) for routing
- Uses [chi](github.com/go-chi/cors) as its CORS
- Uses MOngoDB as its database 