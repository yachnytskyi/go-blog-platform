<div align="left">
  
## Golang Mongo gRPC     
Source code for  Golang Mongo gRPC App.

The project uses:  
**Golang    
Gin  
MongoDB  
gRPC**

Hexagonal architecture, manual dependency injection and abstract factory are implemented in the project.
We have a script in Makefile that allows you to launch the project.

## Initializing
`config/yaml/v1/dev.application.example.yaml` file is provided you as an example of your own config settings,  you need to put them in your `config/yaml/v1/dev.application.yaml` file (you should create it on your own). 

## From the project root directory, run:  
```make up``` 
The API will then be available at  **http://localhost:8080/api/posts  http://localhost:8080/api/users**
You can also find all possible API requests/urls when you launch the project in your server terminal. 

If you need to make rebuild, you can use these commands:  
```make build``` if you prefer a shortcut command from Makefile.   
```docker-compose build``` if you you prefer to enter a full command on your own.
  
After that repeat this command:  
```make up```

## Run server
To run this code, you will need docker and docker-compose installed on your machine. From the root project directory, run:  
```make up``` or   
```make reflex``` use reflex hot reload launch mode  
```make run``` use default launch mode

## Stop Docker Compose services 
```make down``` if you prefer a shortcut command from Makefile.  
```docker-compose down``` if you you prefer to enter a full command on your own.
 
## Ways of possible improvements
I am open for new ideas. Fistful, add unit and integration tests. Secondly, refactoring of the system.

</div>
