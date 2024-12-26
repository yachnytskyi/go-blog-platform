// Connect to the database.
db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE);

// Table names.
const UsersTable = "users";
const PostsTable = "posts";

// Number of users
const usersNumber = 1000; // Number of test users.

// Precomputed bcrypt hashed password
const hashedPassword = "$2b$12$2TnOFrMkPCzvGQZ.PjwdHenVH7kiReJKvXTkp0LPeAC9DnNt8m3ze"; // "somepassword".

// Insert Users first.
let users = [];
for (let i = 0; i < usersNumber; i++) {
    const userId = new ObjectId(); // Generate a new user ID
    users.push({
        _id: userId, // Assign generated ID to user
        username: `test${i}`,
        email: `test${i}@gmail.com`,
        password: hashedPassword, // Use the precomputed hashed password.
        role: "user",
        verified: true,
        created_at: new Date(),
        updated_at: new Date()
    });
}

// Bulk insert users into the users collection.
db[UsersTable].insertMany(users);

// Insert Posts for each user with a random number of posts.
for (let i = 0; i < usersNumber; i++) {
    const userId = users[i]._id; // Retrieve the userId from the previously inserted users
    
    // Generate a random number of posts for the current user (from 0 to 300).
    const postsNumber = Math.floor(Math.random() * 300); // Generates a number between 0 and 300.

    // Insert a random number of posts for the user
    for (let j = 0; j < postsNumber; j++) {
        const postData = {
            _id: new ObjectId(),
            user_id: userId, // Link post to the correct user
            title: `test post${i * postsNumber + j}`,
            content: "Sample content for post",
            image: "https://via.placeholder.com/150",
            username: `test${i}`,
            created_at: new Date(),
            updated_at: new Date()
        };

        // Insert the post into the posts collection.
        db[PostsTable].insertOne(postData);
    }
}
