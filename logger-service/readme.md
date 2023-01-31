The Logger Service 
- The Logger service receives and accepts the data from the authentication, listener and mailer service through the broker service and stores the data in a  mongoDb database


- Uses  [chi router](https://github.com/go-chi/chi/v5) for routing
- Uses [chi](github.com/go-chi/cors) as its CORS
- Uses MOngoDB as its database 