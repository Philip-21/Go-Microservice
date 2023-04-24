# Test Microservice

#### An End-to-End Distributed System that uses different protocols rpc, GRPC , REST API's and Queueing Systems in communicating between services.
#### The main goal of this project is to provide a clear and concise understanding of the communication mechanisms utilized by this microservice architecture-based distributed system. By doing so, users can gain a comprehensive understanding of how the system operates and how these communication mechanisms enable seamless integration between different software services, leading to a more efficient and scalable application..
##### Services include 
##### These are the core services that perform distributed actions 
- [authentication-service](https://github.com/Philip-21/Go-Microservice/tree/master/authentication-service) for Authenticating users
- [Broker Service](https://github.com/Philip-21/Go-Microservice/tree/master/broker-service) is the central node point for handling each request from the client and rendering a response to the client. This Service calls the right service(authentication, listener, logger and mailer )  when a request is called from the front-end  
- [logger-service](https://github.com/Philip-21/Go-Microservice/tree/master/logger-service) receives and accepts the data from the authentication, listener and mailer service ,when each service has been called through the broker service . The data from each service is stored in a  mongoDb database.It also handles the gRPC and RPC actions when called b the broker service. 
- [mail-service](https://github.com/Philip-21/Go-Microservice/tree/master/mail-service) Handles sending of mails
- [listener-service](https://github.com/Philip-21/Go-Microservice/tree/master/listener-service) This service handles queues 

##### other sevices include
- [front-end](https://github.com/Philip-21/Go-Microservice/tree/master/front-end) just displays outputs for each action performed internally
- [project](https://github.com/Philip-21/Go-Microservice/tree/master/project) contains the docker and k8 file to run the application locally