create database food_ai;
use food_ai;

USE mysql;

CREATE USER 'app_user'@'%' IDENTIFIED BY '';

-- Grant all privileges on the specific database to the user
GRANT ALL PRIVILEGES ON food_ai.* TO 'app_user'@'%';

-- Apply changes
FLUSH PRIVILEGES;