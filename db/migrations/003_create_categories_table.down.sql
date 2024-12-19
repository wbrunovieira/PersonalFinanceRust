ALTER TABLE transactions DROP CONSTRAINT fk_category;

ALTER TABLE transactions DROP COLUMN category_id;

ALTER TABLE transactions ADD COLUMN category VARCHAR(50) NOT NULL;

DROP TABLE categories;
