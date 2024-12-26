// Switch to the database.
db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE);

// Create a database user with the required roles.
db.createUser({
  user: process.env.MONGO_INITDB_ROOT_USERNAME,
  pwd: process.env.MONGO_INITDB_ROOT_PASSWORD,
  roles: [
    { role: "readWrite", db: process.env.MONGO_INITDB_DATABASE},
    { role: "dbAdmin", db: process.env.MONGO_INITDB_DATABASE}
  ]
});

// Insert the admin user into the database.
db["users"].insertOne({
  _id: new ObjectId(), // Generate a unique ID for the admin user.
  username: "yachnytskyi",
  email: process.env.ADMIN_EMAIL, 
  password: process.env.ADMIN_HASHED_PASSWORD,
  role: "admin", 
  verified: true,
  created_at: new Date(),
  updated_at: new Date()
});
