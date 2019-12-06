DROP TABLE IF EXISTS user_review;
DROP TABLE IF EXISTS bid_session_log;
DROP TABLE IF EXISTS user_payment_info;
DROP TABLE IF EXISTS bid_session;
DROP TABLE IF EXISTS user_wishlist;
DROP TABLE IF EXISTS user_common;
DROP TABLE IF EXISTS item_image;
DROP TABLE IF EXISTS item;
DROP TABLE IF EXISTS categories;

CREATE TABLE user_common (
    user_id             VARCHAR(255) NOT NULL PRIMARY KEY,
    user_password       VARCHAR(255) NOT NULL,

    user_name           VARCHAR(100),
    user_phone          VARCHAR(15),
    user_email          VARCHAR(255) NOT NULL,
    user_gender         CHAR,
    user_address        VARCHAR(255),
    user_avatar         TEXT,

    user_access_level   int DEFAULT 1,
    user_createAT       DATETIME DEFAULT NOW()
);

CREATE TABLE categories (
    categories_id   INT AUTO_INCREMENT PRIMARY KEY,
    categories_name VARCHAR(255) NOT NULL
);

CREATE TABLE item (
    item_id             BIGINT PRIMARY KEY AUTO_INCREMENT,
    categories_id       INT,
    item_name           VARCHAR(255),
    item_description    TEXT,
    item_condition      VARCHAR(50),
    item_createAt       DATETIME DEFAULT NOW(),
    
    FOREIGN KEY (categories_id) REFERENCES categories(categories_id)
);

CREATE TABLE item_image (
    item_id BIGINT NOT NULL,
    images TEXT,

    FOREIGN KEY (item_id) REFERENCES item(item_id)
);

CREATE TABLE user_payment_info (
    user_id           VARCHAR(255) NOT NULL,
    user_payment_info VARCHAR(100),

    FOREIGN KEY (user_id) REFERENCES user_common(user_id)
);

CREATE TABLE user_wishlist (
    user_id     VARCHAR(255) NOT NULL,
    item_id     BIGINT NOT NULL,
    add_date    DATETIME,

    FOREIGN KEY (user_id) REFERENCES user_common(user_id),
    FOREIGN KEY (item_id) REFERENCES item(item_id) 
);

CREATE TABLE bid_session (
    session_id BIGINT PRIMARY KEY  AUTO_INCREMENT,
    item_id BIGINT NOT NULL,
    seller_id VARCHAR(255) NOT NULL,
    session_start_date DATETIME DEFAULT NOW(),
    session_end_date DATETIME,
    userview_count INT DEFAULT 0,
    winner_id VARCHAR(255),
    minimum_increase_bid INT DEFAULT 0,
    current_bid BIGINT DEFAULT 0,

    FOREIGN KEY (item_id) REFERENCES item(item_id),
    FOREIGN KEY (seller_id) REFERENCES user_common(user_id),
    FOREIGN KEY (winner_id) REFERENCES user_common(user_id)
);

CREATE TABLE bid_session_log (
    user_id VARCHAR(255) NOT NULL,
    session_id BIGINT NOT NULL,
    bid_amount BIGINT NOT NULL,
    bid_date DATETIME,

    FOREIGN KEY (session_id) REFERENCES bid_session(session_id),
    FOREIGN KEY (user_id) REFERENCES user_common(user_id)
);

CREATE TABLE user_review (
    user_writer VARCHAR(255),
    user_target VARCHAR(255),
    session_id  BIGINT,
    review_content TEXT,
    review_score INT(1) NOT NULL,

    FOREIGN KEY (user_writer) REFERENCES user_common(user_id),
    FOREIGN KEY (user_target) REFERENCES user_common(user_id),
    FOREIGN KEY (session_id) REFERENCES bid_session(session_id)
);

INSERT INTO categories VALUES (null,'Video games');
INSERT INTO categories VALUES (null,'Electronics');
INSERT INTO categories VALUES (null,'Computers');
INSERT INTO categories VALUES (null,'Books');
INSERT INTO categories VALUES (null,'Fashions');
INSERT INTO categories VALUES (null,'Luggages');

INSERT INTO item VALUES (null,
1,
'CODE VEIN Steam Keys',
'Key game steam, dùng được trên toàn thế giới',
'NEW',
null);
INSERT INTO item VALUES (null,
1,
'PlayStation 4 Slim 1TB Console ',
'Incredible games; Endless entertainment
All new lighter slimmer PS4
1TB hard drive
All the greatest, games, TV, music and more',
'NEW',
null);
INSERT INTO item VALUES (null,
2,
'Máy hút bụi Philips FC6168',
'Thiết kế nhỏ gọn, tiện lợi

Máy Hút Bụi Cầm Tay Philips FC6168 là vật dụng có khả năng hút bụi và lau sàn chỉ trong một lần di chuyển với thiết kế 2 trong 1 tiện lợi và tay cầm thuận tiện giúp bạn nhanh chóng lau dọn vết bẩn mỗi ngày. Thích hợp sử dụng trên tất cả các bề mặt sàn, như sàn gỗ, thảm, Hoặc bạn có thể lắp ngăn chứa nước để lau chùi các bề mặt sàn cứng.',
'NEW',
null);
INSERT INTO item VALUES (null,
2,
'Pin sạc dự phòng Li-ion 26800mAh Anker PowerCore+ A1375 Đen',
null,
'USED',
null);

INSERT INTO item_image VALUES(1,'/view/images/codevein1');
INSERT INTO item_image VALUES(1,'/view/images/codevein2');
INSERT INTO item_image VALUES(1,'/view/images/codevein3');
INSERT INTO item_image VALUES(2,'/view/images/playstation4_1');
INSERT INTO item_image VALUES(2,'/view/images/playstation4_2');
INSERT INTO item_image VALUES(2,'/view/images/playstation4_3');
INSERT INTO item_image VALUES(3,'/view/images/may-hut-bui-philips-fc6168-800-1');
INSERT INTO item_image VALUES(3,'/view/images/may-hut-bui-philips-fc6168-800-2');
INSERT INTO item_image VALUES(3,'/view/images/may-hut-bui-philips-fc6168-800-3');
INSERT INTO item_image VALUES(4,'/view/images/sacLi-ion');

INSERT INTO user_common VALUES("tester66",'6666','Death Click','6666666666','death@click.hell','M','hell',null,2,null);
INSERT INTO user_common VALUES("tester67",'6666','Death Click CLone','6666666666','death@click.hell','M','hell',null,2,null);
INSERT INTO user_common VALUES("tester68",'6666','Death Click CLone','6666666666','death@click.hell','M','hell',null,2,null);


INSERT INTO bid_session VALUES(null,1,"tester66",'2019-11-22 18:00:04','2020-11-22 18:00:04',null,null,null,null);
INSERT INTO bid_session VALUES(null,2,"tester66",'2019-11-22 18:00:04','2020-11-22 18:00:04',null,null,null,null);
INSERT INTO bid_session VALUES(null,3,"tester66",'2019-11-22 18:00:04','2020-11-22 18:00:04',null,null,null,null);
INSERT INTO bid_session VALUES(null,4,"tester66",'2019-11-22 18:00:04','2020-11-22 18:00:04',null,null,null,null);


INSERT INTO bid_session_log VALUES("tester67",1,330000,'2019-11-22 18:12:22');
INSERT INTO bid_session_log VALUES("tester67",2,1500000,'2019-11-22 18:01:10');
INSERT INTO bid_session_log VALUES("tester68",2,1805000,'2019-11-22 19:45:39');
INSERT INTO bid_session_log VALUES("tester67",2,2210000,'2019-11-22 21:04:09');
INSERT INTO bid_session_log VALUES("tester68",2,2800000,'2019-11-23 11:12:40');
INSERT INTO bid_session_log VALUES("tester67",2,800000,'2019-11-22 18:12:12');
INSERT INTO bid_session_log VALUES("tester68",3,1600001,'2019-11-23 17:17:17');
