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
    phone_number VARCHAR(13) UNIQUE NOT NULL,
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
31,MOY UNIVERSAL 20/50,900GR,14,50000,15000,30%,0,1311,0,2.622%
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    size INT NOT NULL, 
    type product_type DEFAULT 'countable' NOT NULL, 
    price FLOAT NOT NULL,
    selling_price FLOAT,
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







-- qarzni hisoblash
CREATE OR REPLACE FUNCTION recalc_user_debt()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE users u
    SET debt = COALESCE(took_sum, 0)
    FROM (
        SELECT
            user_id,
            SUM(CASE WHEN debt_type = 'took' AND deleted_at = 0 THEN amount ELSE 0 END) AS took_sum
        FROM debt_logs
        WHERE deleted_at = 0
        GROUP BY user_id
    ) d
    WHERE u.id = d.user_id
      AND u.id = COALESCE(NEW.user_id, OLD.user_id);

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS debt_logs_after_change ON debt_logs;

CREATE TRIGGER debt_logs_after_change
AFTER INSERT OR UPDATE OR DELETE ON debt_logs
FOR EACH ROW
EXECUTE FUNCTION recalc_user_debt();



-- total_caount
CREATE OR REPLACE FUNCTION recalc_bucket_total(p_bucket_id UUID)
RETURNS void AS $$
DECLARE
    v_total FLOAT;
BEGIN
    SELECT COALESCE(SUM(count * price), 0)
    INTO v_total
    FROM bucket_item
    WHERE bucket_id = p_bucket_id
      AND deleted_at = 0;

    UPDATE buckets
    SET total_price = v_total,
        updated_at = NOW()
    WHERE id = p_bucket_id;
END;
$$ LANGUAGE plpgsql;




CREATE OR REPLACE FUNCTION bucket_item_after_change()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        PERFORM recalc_bucket_total(OLD.bucket_id);
    ELSE
        PERFORM recalc_bucket_total(NEW.bucket_id);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_bucket_item_change
AFTER INSERT OR UPDATE OR DELETE
ON bucket_item
FOR EACH ROW
EXECUTE FUNCTION bucket_item_after_change();
