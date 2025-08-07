-- Удаляем таблицы, если существуют (в правильном порядке)
-- DROP TABLE IF EXISTS players_assets;
DROP TABLE IF EXISTS assets CASCADE;
DROP TABLE IF EXISTS assets_type CASCADE;

-- Создаем таблицу типов активов
CREATE TABLE assets_type (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL UNIQUE
);

INSERT INTO assets_type (id, type) VALUES
(1, 'small_deal'),
(2, 'big_deal');

-- Создаем таблицу assets с внешним ключом
CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    title VARCHAR(50) NOT NULL,
    descr VARCHAR(100) NOT NULL,
    type_id INTEGER,
    price INTEGER DEFAULT 0,
    cashflow INTEGER DEFAULT 0,
    CONSTRAINT fk_asset_type FOREIGN KEY (type_id) REFERENCES assets_type(id) ON DELETE RESTRICT
);

create table players_assets (
    asset_id integer not null,
    player_id integer not null,
    primary key (asset_id, player_id),
    constraint fk_asset foreign key (asset_id) references assets(id) on delete cascade
);

-- Создаем индекс для type_id
CREATE INDEX idx_assets_type_id ON assets(type_id);

-- alter table cashflow.players add constraint fk_character_id foreign key (character_id) references characters(id) on delete cascade;
-- alter table cashflow.players add constraint fk_character_id foreign key (character_id) references characters(id) on delete cascade;