CREATE DATABASE IF NOT EXISTS general;

USE general;

DROP TABLE IF EXISTS t_user, t_user_unverified;

CREATE TABLE IF NOT EXISTS t_user (
    id INTEGER AUTO_INCREMENT primary key,
    user_name VARCHAR(20) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    phone VARCHAR(11),
    gender VARCHAR(1),
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    index idx_email (email),
    index idx_email_activated (gender),
    index idx_user_name (user_name)
);

CREATE TABLE IF NOT EXISTS t_user_unverified (
    id INTEGER AUTO_INCREMENT primary key,
    user_name VARCHAR(20) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    index idx_email (email, created_at)
);