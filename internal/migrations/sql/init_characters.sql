DROP TABLE IF EXISTS characters CASCADE;

-- Создаем таблицу типов активов
CREATE TABLE characters (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    job VARCHAR(50) NOT NULL,
    salary INTEGER NOT NULL DEFAULT 0,
    taxes INTEGER NOT NULL DEFAULT 0,
    child_expenses INTEGER NOT NULL DEFAULT 0,
    other_expenses INTEGER NOT NULL DEFAULT 0
    
);

-- alter table players add constraint fk_character_id foreign key (character_id) references characters(id) on delete cascade;
-- alter table players add constraint fk_character_id foreign key (character_id) references characters(id) on delete cascade;
-- alter table players add constraint fk_character_id foreign key (character_id) references characters(id) on delete cascade;