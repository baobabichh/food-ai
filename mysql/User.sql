CREATE TABLE Users
(
    ID BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    Email VARCHAR(100) UNIQUE,
    Password VARCHAR(100),
    UUID VARCHAR(40) UNIQUE,
    CreateTS TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);