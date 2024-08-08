CREATE TABLE investors (
    id VARCHAR PRIMARY KEY,
    name VARCHAR
);

CREATE TABLE deposits (
    id VARCHAR PRIMARY KEY,
    investor_id VARCHAR,
    FOREIGN KEY (investor_id) REFERENCES investors(id)
);

CREATE TABLE pots (
    id VARCHAR PRIMARY KEY,
    deposit_id VARCHAR,
    name VARCHAR,
    FOREIGN KEY (deposit_id) REFERENCES deposits(id)
);

CREATE TABLE accounts (
    id VARCHAR PRIMARY KEY,
    pot_id VARCHAR,
    wrapper_type INTEGER,
    nominal_amount INTEGER,
    FOREIGN KEY (pot_id) REFERENCES pots(id)
);

CREATE TABLE receipts (
    id VARCHAR PRIMARY KEY,
    account_id VARCHAR,
    amount INTEGER,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);
