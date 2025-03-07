CREATE TABLE FoodRecognitionRequests
(
    ID BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ImgBase64 LONGTEXT,
    ImgType VARCHAR(20),
    Status TINYINT DEFAULT 0, -- 1 pending, 2 processing, 3 error, 4 success
    CreateTS TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);