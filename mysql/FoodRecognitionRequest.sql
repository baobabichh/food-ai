CREATE TABLE FoodRecognitionRequests
(
    ID BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    ImgBase64 LONGTEXT,
    ImgType VARCHAR(20),
    Status TINYINT DEFAULT 0, -- 1 pending, 2 processing, 3 error, 4 success
    Response LONGTEXT,
    CreateTS TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


alter table FoodRecognitionRequests add column UserID BIGINT UNSIGNED not null;
alter table FoodRecognitionRequests add constraint fk_FoodRecognitionRequests_UserID foreign key(UserID) references Users(ID);