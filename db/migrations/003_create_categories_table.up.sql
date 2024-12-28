CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
    type VARCHAR(10) NOT NULL 
);



INSERT INTO categories (name, type) VALUES 
    ('Alimentação', 'expense'),
    ('Salário', 'income'),
    ('Saúde', 'expense'),
    ('Investimentos', 'income');

ALTER TABLE transactions 
    ADD COLUMN category_id SET NOT NULL;;

ALTER TABLE transactions 
    DROP COLUMN category;

ALTER TABLE transactions 
    ADD CONSTRAINT fk_category
    FOREIGN KEY (category_id) 
    REFERENCES categories(id) 
    ON DELETE CASCADE;
