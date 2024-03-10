DROP DATABASE IF EXISTS apirest;
CREATE DATABASE apirest;

USE apirest;

CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

INSERT INTO user VALUES (1, "julian777", "12345", "julian@mail.com");
INSERT INTO user VALUES (2, "burgos777", "12345", "burgos@mail.com");
INSERT INTO user VALUES (3, "alejo777", "12345", "alejo@mail.com");