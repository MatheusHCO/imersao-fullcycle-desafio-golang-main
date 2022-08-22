CREATE TABLE accounts (
    id uuid NOT NULL,
    number VARCHAR NOT NULL UNIQUE,
    balance float NOT NULL,
    PRIMARY KEY (id)
);
CREATE TABLE transactions (
    id uuid NOT NULL,
    `from` varchar NOT NULL,
    `to` varchar NOT NULL,
    amount float NOT NULL,
    PRIMARY KEY (id)
);