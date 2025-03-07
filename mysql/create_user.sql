create database food_ai;
use food_ai;

USE mysql;

CREATE USER 'user'@'%' IDENTIFIED BY 'pass';

-- Grant all privileges on the specific database to the user
GRANT ALL PRIVILEGES ON food_ai.* TO 'user'@'%';

-- Apply changes
FLUSH PRIVILEGES;