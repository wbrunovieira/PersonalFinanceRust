CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO categories (name) VALUES 
    ('Alimentacao'),
    ('Saude'),
    ('Transporte'),
    ('Lazer'),
    ('Educacao'),
    ('Moradia'),
    ('Investimento'),
    ('Outros');

ALTER TABLE transactions 
    ADD COLUMN category_id INT NOT NULL;

ALTER TABLE transactions 
    DROP COLUMN category;

ALTER TABLE transactions 
    ADD CONSTRAINT fk_category
    FOREIGN KEY (category_id) 
    REFERENCES categories(id) 
    ON DELETE CASCADE;
