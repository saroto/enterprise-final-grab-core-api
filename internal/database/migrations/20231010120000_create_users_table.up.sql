DROP TABLE IF EXISTS public.users;

CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    vehicleName VARCHAR(255) DEFAULT NULL,
    vehicleNumber VARCHAR(255) DEFAULT NULL,
    role VARCHAR(255) NOT NULL DEFAULT 'user',
    password TEXT NOT NULL
);

CREATE TABLE public.tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    debit DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    credit DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    CONSTRAINT unique_user_id UNIQUE (user_id)
);

CREATE TABLE public.transactions (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES public.accounts(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL,
    method VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

