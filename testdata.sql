DROP TABLE IF EXISTS bid_session_log;
DROP TABLE IF EXISTS user_payment_info;
DROP TABLE IF EXISTS bid_session;
DROP TABLE IF EXISTS user_wishlist;
DROP TABLE IF EXISTS user_common;
DROP TABLE IF EXISTS item_image;
DROP TABLE IF EXISTS item;
DROP TABLE IF EXISTS categories;

CREATE TABLE user_common (
    user_id             int PRIMARY KEY AUTO_INCREMENT,
    user_name           VARCHAR(100),
    user_phone          VARCHAR(15),
    user_birth          DATE,
    user_gender         CHAR,
    user_adress         VARCHAR(255),
    user_loginID        VARCHAR(32) NOT NULL,
    user_password       VARCHAR(255) NOT NULL,
    user_access_level   int NOT NULL,
    user_login_token    TEXT
);

CREATE TABLE categories (
    categories_id   int NOT NULL PRIMARY KEY,
    categories_name VARCHAR(255) NOT NULL
);

CREATE TABLE item (
    item_id             int PRIMARY KEY AUTO_INCREMENT,
    categories_id       int,
    item_name           VARCHAR(255),
    item_description    TEXT,
    item_condition      VARCHAR(50),
    item_sale_status    VARCHAR(30),
    item_add_time       DATE,
    
    FOREIGN KEY (categories_id) REFERENCES categories(categories_id)
);

CREATE TABLE item_image (
    item_id int,
    image_link TEXT,

    FOREIGN KEY (item_id) REFERENCES item(item_id)
);

CREATE TABLE user_payment_info (
    user_id           int,
    user_payment_info VARCHAR(100),

    FOREIGN KEY (user_id) REFERENCES user_common(user_id)
);

CREATE TABLE user_wishlist (
    user_id int,
    item_id int,

    FOREIGN KEY (user_id) REFERENCES user_common(user_id),
    FOREIGN KEY (item_id) REFERENCES item(item_id) 
);

CREATE TABLE bid_session (
    session_id int PRIMARY KEY  AUTO_INCREMENT,
    item_id int,
    seller_id int,
    session_status VARCHAR(30),
    session_start_date DATE,
    session_end_date DATE,

    FOREIGN KEY (item_id) REFERENCES item(item_id),
    FOREIGN KEY (seller_id) REFERENCES user_common(user_id)
);

CREATE TABLE bid_session_log (
    session_id INT,
    user_id INT,
    bid_amount FLOAT(14,2),
    bid_date DATETIME,

    FOREIGN KEY (session_id) REFERENCES bid_session(session_id),
    FOREIGN KEY (user_id) REFERENCES user_common(user_id)
);

INSERT INTO categories VALUES (1,'Video games');
INSERT INTO categories VALUES (2,'Electronics');
INSERT INTO categories VALUES (3,'Computers');
INSERT INTO categories VALUES (4,'Books');
INSERT INTO categories VALUES (5,'Fashions');
INSERT INTO categories VALUES (6,'Luggages');

INSERT INTO item VALUES (1,1,'CODE VEIN Steam Keys','Key game steam, dùng được trên toàn thế giới','NEW','ĐANG ĐẤU GIÁ','2019-11-22');
INSERT INTO item VALUES (2,1,'PlayStation 4 Slim 1TB Console ','Incredible games; Endless entertainment
All new lighter slimmer PS4
1TB hard drive
All the greatest, games, TV, music and more','NEW','ĐANG ĐẤU GIÁ','2019-11-25');
INSERT INTO item VALUES (3,2,'Máy hút bụi Philips FC6168','Thiết kế nhỏ gọn, tiện lợi

Máy Hút Bụi Cầm Tay Philips FC6168 là vật dụng có khả năng hút bụi và lau sàn chỉ trong một lần di chuyển với thiết kế 2 trong 1 tiện lợi và tay cầm thuận tiện giúp bạn nhanh chóng lau dọn vết bẩn mỗi ngày. Thích hợp sử dụng trên tất cả các bề mặt sàn, như sàn gỗ, thảm, Hoặc bạn có thể lắp ngăn chứa nước để lau chùi các bề mặt sàn cứng.','NEW','ĐANG ĐẤU GIÁ','2019-11-23');
INSERT INTO item VALUES (4,2,'Pin sạc dự phòng Li-ion 26800mAh Anker PowerCore+ A1375 Đen',null,'USED','ĐÃ BÁN','2019-11-20');

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

INSERT INTO user_common VALUES(1,'Trương Quang Hiếu','0123456789','2002-02-04','M',null,'tester01','123',3,null);
INSERT INTO user_common VALUES(2,'Trần Ngọc Quý','1456238900','1998-02-20','M',null,'tester02','123',3,null);
INSERT INTO user_common VALUES(3,'Ricardo Milos','0452854491','1985-05-05','M',null,'tester03','123',1,null);
INSERT INTO user_common VALUES(4,'Death Click','1856040012','2001-12-12','M',null,'tester04','123',1,null);

INSERT INTO user_wishlist VALUES(1,1);
INSERT INTO user_wishlist VALUES(1,2);
INSERT INTO user_wishlist VALUES(2,2);
INSERT INTO user_wishlist VALUES(3,3);
INSERT INTO user_wishlist VALUES(3,4);

INSERT INTO bid_session VALUES(1,1,4,'CURRENTLY RUNNING',null,null);
INSERT INTO bid_session VALUES(2,2,4,'CURRENTLY RUNNING',null,null);
INSERT INTO bid_session VALUES(3,3,4,'CURRENTLY RUNNING',null,null);

INSERT INTO bid_session_log VALUES(1,1,330000,'2019-11-22 18:12:22');
INSERT INTO bid_session_log VALUES(2,1,1500000,'2019-11-22 18:01:10');
INSERT INTO bid_session_log VALUES(2,2,1805000,'2019-11-22 19:45:39');
INSERT INTO bid_session_log VALUES(2,1,2210000,'2019-11-22 21:04:09');
INSERT INTO bid_session_log VALUES(2,2,2800000,'2019-11-23 11:12:40');
INSERT INTO bid_session_log VALUES(3,2,800000,'2019-11-22 18:12:12');
INSERT INTO bid_session_log VALUES(3,3,1600001,'2019-11-23 17:17:17');
