CREATE TYPE role AS ENUM('admin', 'user');

CREATE TYPE product_type AS ENUM('g', 'ml', 'countable');

CREATE TYPE payment_type AS ENUM('cash', 'card');

CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled', 'returned');

CREATE TYPE debt_type AS ENUM ('took', 'gave');

CREATE EXTENSION postgis;


CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    password TEXT,
    phone_number VARCHAR(20) NOT NULL,
    debt BIGINT NOT NULL DEFAULT 0,
    role role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE debt_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    reason TEXT NOT NULL,
    debt_type debt_type NOT NULL DEFAULT 'took',
    time_taken TIMESTAMP NOT NULL DEFAULT NOW(),
    given_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    size INT NOT NULL, 
    type product_type DEFAULT 'countable' NOT NULL, 
    price FLOAT NOT NULL,
    img_url TEXT NOT NULL,
    count INT NOT NULL,
    sales_count INT NOT NULL DEFAULT 0,
    category_id UUID NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS buckets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_price FLOAT NOT NULL DEFAULT 0,
    status BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS bucket_item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bucket_id UUID NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    count INT NOT NULL,
    price FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bucket_id UUID NOT NULL REFERENCES buckets(id) ON DELETE CASCADE,
    status order_status NOT NULL DEFAULT 'pending',
    location GEOGRAPHY(POINT, 4326),
    description TEXT NOT NULL,
    payment_type payment_type DEFAULT 'cash' NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);

