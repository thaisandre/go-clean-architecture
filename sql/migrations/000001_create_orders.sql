CREATE TABLE IF NOT EXISTS orders (
    id          VARCHAR(255) NOT NULL PRIMARY KEY,
    price       DECIMAL(10, 2) NOT NULL,
    tax         DECIMAL(10, 2) NOT NULL,
    final_price DECIMAL(10, 2) NOT NULL
);

INSERT INTO orders (id, price, tax, final_price) VALUES ('f937b55b-f71a-4ae6-b2cb-3d02752b6c0d', 100.50, 0.50, 101.00)
