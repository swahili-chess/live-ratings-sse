CREATE TABLE static(
    id SERIAL PRIMARY KEY,
    lichess_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    rapid INT NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
);


CREATE TABLE dynamic (
    id SERIAL PRIMARY KEY,
    lichess_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    rapid INT NOT NULL,
    created_at TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP(0) NOT NULL DEFAULT NOW()
    
);