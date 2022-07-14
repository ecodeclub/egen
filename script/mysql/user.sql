DROP DATABASE if exists user;

CREATE DATABASE `user` charset=utf8;

CREATE TABLE `user`.`user_account`(
    `id` INT  AUTO_INCREMENT PRIMARY KEY,
    `password` VARCHAR(15),
    `login` VARCHAR(25),
    `username` VARCHAR(20)
)
