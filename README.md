<div align="left">
  
## Golang Mongo gRPC     
Source code for  Golang Mongo gRPC App.

The project uses:  
**Golang    
Gin  
MongoDB  
gRPC**

Hexagonal architecture, manual dependency injection and abstract factory are implemented in the project.  

## Initializing
`config/yaml/v1/local.dev.application.example.yaml` and `config/yaml/v1/docker.dev.application.example.yaml` file are provided you as the examples of your own config settings, you need to put them in your `config/yaml/v1/local.dev.application.example.yaml` and/or 
`config/yaml/v1/docker.dev.application.example.yaml` (you should create them on your own). 

## From the project root directory, run:  
```make up``` For docker environment.  ```make local ``` For local environment     
    
The API will then be available at  **http://localhost:8080/api/posts  http://localhost:8080/api/users**  
You can also find all possible API requests/urls when you launch the project in your server terminal. 

If you need to make rebuild, you can use these commands:  
```make build``` if you prefer a shortcut command from Makefile.   
```docker-compose build``` if you you prefer to enter a full command on your own.
  
After that repeat this command:
```make up``` For docker environment.  ```make local ``` For local environment    

## Run server
To run this code, you will need docker and docker-compose installed on your machine. From the root project directory, run:  
```make up``` use default dockerized launch mode     
```make local``` use local launch mode 

## Stop Docker Compose services 
```make down``` if you prefer a shortcut command from Makefile.  
```docker-compose down``` if you you prefer to enter a full command on your own.
 
## Ways of possible improvements
I am open for new ideas. At first add unit and integration tests. Secondly, refactoring of the system.

</div>
