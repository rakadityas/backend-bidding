DROP TABLE IF EXISTS product;
CREATE TABLE IF NOT EXISTS product (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  product_name varchar(250) NOT NULL,
  image_url varchar(250) NOT NULL,
  "status" INT NOT NULL,
  create_time TIMESTAMP,
  update_time TIMESTAMP
);

DROP TABLE IF EXISTS "user";
CREATE TABLE IF NOT EXISTS "user" (
  id SERIAL PRIMARY KEY,
  user_type INT NOT NULL,
  username VARCHAR(20) NOT NULL UNIQUE,
  "password" VARCHAR(30) NOT NULL UNIQUE,
  "status" INT,
  balance INT,
  create_time TIMESTAMP,
  update_time TIMESTAMP
);

DROP TABLE IF EXISTS timewindow;
CREATE TABLE IF NOT EXISTS timewindow (
  id SERIAL PRIMARY KEY,
  auction_id INT,
  start_time TIMESTAMP,
  end_time TIMESTAMP,
  "status" INT,
  create_time TIMESTAMP,
  update_time TIMESTAMP
);


DROP TABLE IF EXISTS bid_collection;
CREATE TABLE IF NOT EXISTS bid_collection (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  auction_id INT NOT NULL,
  current_bid INT NOT NULL,
  payment_id INT NOT NULL,
  create_time TIMESTAMP,
  update_time TIMESTAMP
);

DROP TABLE IF EXISTS auction;
CREATE TABLE IF NOT EXISTS auction (
  id SERIAL PRIMARY KEY,
  product_id INT NOT NULL, 
  winner_user_id INT,
  multiplier INT NOT NULL,
  "status" INT NOT NULL,
  create_time TIMESTAMP,
  update_time TIMESTAMP
);

DROP TABLE IF EXISTS payment;
CREATE TABLE IF NOT EXISTS payment (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL, 
  amount INT,
  "status" INT NOT NULL,
  create_time TIMESTAMP,
  update_time TIMESTAMP
);

INSERT INTO product(user_id, product_name, image_url, "status", create_time, update_time) VALUES 
(1, 'Dummy product 1', 'https://i.imgur.com/ibRuYT4.jpeg', 1, 'NOW()', 'NOW()'),
(2, 'Dummy product 2', 'https://i.imgur.com/HcGL77B.jpeg', 1, 'NOW()', 'NOW()');


INSERT INTO "user" (user_type, username, "password", "status", create_time, balance) VALUES
(1, 'testuser1', 'testuser1', 1, now(), 100000),
(0, 'testuser2', 'testuser2', 1, now(), 100000);

INSERT INTO auction (product_id, winner_user_id, multiplier, "status", create_time) VALUES
(1, NULL, 10000, 1, now()),
(2, NULL, 10000, 1, now());

INSERT INTO timewindow (auction_id, start_time, end_time, "status", create_time, update_time) VALUES
(1, '2022-02-23 11:01:34.581', '2022-03-19 11:01:34.581', 1, 'NOW()', 'NOW()'),
(2, '2022-02-23 11:01:34.581', '2022-03-19 11:01:34.581', 1, 'NOW()', 'NOW()');
