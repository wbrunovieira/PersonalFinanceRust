CREATE TABLE IF NOT EXISTS  categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    category_type VARCHAR(10) NOT NULL
);

INSERT INTO categories (name, category_type) VALUES 
    ('Alimentação', 'expense'),
    ('Salário', 'income'),
    ('Saúde', 'expense'),
    ('Investimentos', 'income'),
    ('Moradia', 'expense'),
    ('Lazer', 'expense'),
    ('Transporte', 'expense'),
    ('Educação', 'expense'),
    ('Vestuário', 'expense'),
    ('Poupança', 'income'),
    ('Tesouro Direto', 'investment'),
    ('Ações', 'investment'),
    ('FII', 'investment'),
    ('Renda Fixa', 'investment'),
    ('Previdência Privada', 'investment');
