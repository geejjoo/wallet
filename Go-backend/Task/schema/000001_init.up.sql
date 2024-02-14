CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE WALLET
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    balance NUMERIC(8,3) NOT NULL
);
CREATE TABLE HISTORY
(
    id serial PRIMARY KEY,
    from_id uuid NOT NULL,
    to_id uuid NOT NULL,
    amount NUMERIC(8,3) not null,
    time timestamp
);
