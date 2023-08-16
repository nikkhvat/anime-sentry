CREATE ROLE postgres WITH LOGIN PASSWORD 'postgres';

CREATE DATABASE anime_schedule WITH OWNER = postgres;
GRANT ALL PRIVILEGES ON DATABASE anime_schedule TO postgres;