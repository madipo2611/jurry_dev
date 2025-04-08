CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100),
                       login VARCHAR(50) NOT NULL UNIQUE,
                       password CHAR(250) NOT NULL,
                       date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       balans INTEGER NOT NULL DEFAULT 0,
                       status VARCHAR(30) NOT NULL DEFAULT 'active',
                       role VARCHAR(10) NOT NULL DEFAULT 'normal',
                       last_seen TIMESTAMP NULL DEFAULT NULL,
                       gender VARCHAR(50) NOT NULL,
                       language VARCHAR(10) NOT NULL DEFAULT 'en',
                       active_status_online BOOLEAN NOT NULL DEFAULT TRUE,
                       posts_privacy SMALLINT NOT NULL DEFAULT 1,
                       allow_dm SMALLINT NOT NULL DEFAULT 1,
                       allow_comments SMALLINT NOT NULL DEFAULT 1
);

CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       userID INTEGER NOT NULL DEFAULT 0,
                       image VARCHAR(150) NOT NULL UNIQUE,
                       text VARCHAR(150) NOT NULL UNIQUE,
                       likes INTEGER NOT NULL DEFAULT 0,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (userID) REFERENCES users(id)
);

CREATE INDEX ON posts (userID);
