DROP TABLE IF EXISTS market CASCADE;

CREATE TABLE market (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    title VARCHAR(50) NOT NULL,
    descr VARCHAR(100) NOT NULL,
    type_id INTEGER,
    sell_cost INTEGER DEFAULT 0,

    CONSTRAINT fk_market_type FOREIGN KEY (type_id) REFERENCES assets_type(id) ON DELETE RESTRICT
);

CREATE INDEX idx_market_type_id ON market(type_id);