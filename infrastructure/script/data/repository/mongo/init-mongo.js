// Switch to the database.
db = db.getSiblingDB('golang_mongodb');

// Create a user with the required roles.
db.createUser({
  user: "root",
  pwd: "root",
  roles: [
    { role: "readWrite", db: "golang_mongodb" }, 
    { role: "dbAdmin", db: "golang_mongodb" }    
  ]
});
