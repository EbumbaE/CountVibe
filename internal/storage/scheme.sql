CREATE TYPE composition AS (
    calories FLOAT
	proteins FLOAT
	fats FLOAT
	carbohydrates FLOAT
);

CREATE TABLE IF NOT EXIST users (
    user_id     INT PRIMARY KEY UNIQUE NOT NULL,
    diary_id    INT PRIMARY KEY UNIQUE NOT NULL,
    username    VARCHAR(30) NOT NULL,
    password    VARCHAR(30) NOT NULL,
);

CREATE TABLE IF NOT EXIST products (
    product_id       INT PRIMARY KEY UNIQUE NOT NULL,
	name             VARCHAR(30) NOT NULL,
	unit_composition composition 
	unit string      VARCHAR(10) NOT NULL,
);

CREATE TABLE IF NOT EXIST diary (
    diary_id    INT PRIMARY KEY UNIQUE NOT NULL,
    date        VARCHAR(10) NOT NULL,
    meal_name   VARCHAR(30) NOT NULL,
    product_id  integer PRIMARY KEY UNIQUE NOT NULL,
    amount      FLOAT NOT NULL,       
);