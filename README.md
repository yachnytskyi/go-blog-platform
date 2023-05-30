<div align="left">
  
## Golang Mongo gRPC     
Source code for  Golang Mongo gRPC App.

The project uses:

**Golang    
Gin  
MongoDB  
gRPC**

We have a script in Makefile that allows you to launch the project.

## Runing the Application

  
`example.env` file is provided you as an example of your own environment variables, which you you need to put in your `app.env` file (you should create it on your own). 


## From the project root director, run:

```make up```

The API will then be available at  **http://localhost:8080/api/posts and http://localhost:8080/api/users**

You can also find all possible API requests/urls when you launch the project in your server terminal. 

If you need to make rebuild, you have to use this command:

```docker-compose build``` 
  
After that repeat ```docker-compose up``` or ```make up``` commands for launching the project.


## gRPC

To launch gRPC server, you need to comment 'Gin server' out and uncomment 'gRPC server' out in `cmd/server/main.go` file. After successful launch,
use this command:

```evans --host localhost --port 8081 -r repl```

## Run

To run this code, you will need docker and docker-compose installed on your machine. In the project root, run:  

```make up```    

```make run```
  
# Ways of possible improvements
I would be grateful for any help you could provide. First of all, I would implement Abstract Factory pattern, to give us an ability to easily switch between repositories and delivery tools. At the current state the settings are hardcoded in `cmd/server/main.go` file, they should be in a config file and defined on first launch of the app. Secondly, add unit and integration tests. Thirdly, fix a bug with launching the application from `Docker Compose`. We have a problem here, that we can't launch the app directly from `Docker Compose`, we need to use `go run main.go` command. I had tried to fix this problem, but haven't had any success. Fourthly, to make a general refactoring of the system.

</div>
