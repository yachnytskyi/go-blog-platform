// Switch to the database.
db = db.getSiblingDB('your database name');

// Create a user with the required roles.
db.createUser({
  user: "your username",
  pwd: "your password",
  roles: [
    { role: "readWrite", db: "your database name" }, 
    { role: "dbAdmin", db: "your database name" }    
  ]
});
