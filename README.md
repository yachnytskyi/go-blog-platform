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
  
To activate email sending features, you need to put your email provider credentials in app.env (you have an example in the end of `app.env` file).

  
`example.env` file is provided you as an example of your own environment variables, which you you need to put in your `app.env` file (you should create it on your own). 


## From the project root director, run:

```make up``` if you prefer a shortcut command from Makefile.

```docker-compose up``` if you you prefer to enter a full command on your own.

The API will then be available at  **http://localhost:8080/api/posts and http://localhost:8080/api/users**

You can also find all possible API requests/urls when you launch the project in your server terminal. 

If you need to make rebuild, you can use these commands:

```make build``` if you prefer a shortcut command from Makefile.

```docker-compose build``` if you you prefer to enter a full command on your own.
  
After that repeat these commands:

```make up``` if you prefer a shortcut command from Makefile.

```docker-compose up``` if you you prefer to enter a full command on your own.


## gRPC server

To launch gRPC server, you need to comment `Gin server` out and uncomment `gRPC server` out in `cmd/server/main.go` file. After a successful launch,
please use this command:

```evans --host localhost --port 8081 -r repl```

If you'd like to return Gin API (or whatever REST API server you'd prefer to use), you should comment `gRPC server` out and uncomment `your server` out (for example our current `Gin server`) in `cmd/server/main.go` file.

## Run server

To run this code, you will need docker and docker-compose installed on your machine. In the project root, run:  

```make up```    

```make reflex``` (if you'd like to use the reflex hot reload launch mode of the server)

```make run``` (if you'd like to use `the default launch mode` of the server)

## Stop server

```make down``` if you prefer a shortcut command from Makefile.

```docker-compose down``` if you you prefer to enter a full command on your own.

  
# Ways of possible improvements
I would be grateful for any help you could provide. First of all, I would implement Abstract Factory pattern, to give us an ability to easily switch between repositories and delivery tools. At the current state the settings are hardcoded in `cmd/server/main.go` file, they should be in the config file and defined on first launch of the app. Secondly, add unit and integration tests. Thirdly, to make a general refactoring of the system.

</div>
