CREATE TYPE states as enum ('created', 'deleted');

ALTER TABLE users 
ALTER COLUMN status TYPE VARCHAR 
USING status::varchar;

ALTER TABLE users 
ALTER COLUMN status TYPE states 
USING status::states;