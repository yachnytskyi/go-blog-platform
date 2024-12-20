// Switch to the database.
db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE || 'default_db');

// Create a user with the required roles.
db.createUser({
  user: process.env.MONGO_INITDB_ROOT_USERNAME || "default_user",
  pwd: process.env.MONGO_INITDB_ROOT_PASSWORD || "default_password",
  roles: [
    { role: "readWrite", db: process.env.MONGO_INITDB_DATABASE || "default_db" },
    { role: "dbAdmin", db: process.env.MONGO_INITDB_DATABASE || "default_db" }
  ]
});
