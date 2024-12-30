CREATE TABLE IF NOT EXISTS projections (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    description TEXT NOT NULL,
    category_id INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_recurring BOOLEAN NOT NULL DEFAULT false,
    end_month DATE,
    date TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);