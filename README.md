# postgres db
CREATE DATABASE jwt;

    CREATE USER jwt WITH SUPERUSER CREATEDB CREATEROLE LOGIN PASSWORD 'jwt';

    ALTER ROLE jwt SET client_encoding TO 'utf8';

    ALTER ROLE jwt SET default_transaction_isolation TO 'read committed';

    ALTER ROLE jwt SET timezone TO 'UTC';

    GRANT ALL PRIVILEGES ON DATABASE jwt TO jwt;