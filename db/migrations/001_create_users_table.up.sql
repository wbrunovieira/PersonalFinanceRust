CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    age INT NOT NULL
);


 INSERT INTO users (id, name, email, age) VALUES (1, 'Bruno', 'bruno@example.com', 50);