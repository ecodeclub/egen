DROP DATABASE if exists user_infor;

CREATE DATABASE `user_infor` charset=utf8;

CREATE TABLE `user_infor`.`user`(
    `id` INT  AUTO_INCREMENT PRIMARY KEY,
    `password` VARCHAR(15),
    `login` VARCHAR(25),
    `username` VARCHAR(20)
);

CREATE TABLE `user_infor`.`user_first`(
    `id` INT  AUTO_INCREMENT PRIMARY KEY,
    `password` VARCHAR(15),
    `login` VARCHAR(25),
    `username` VARCHAR(20)
);

CREATE TABLE `user_infor`.`user_second`(
    `id` INT  AUTO_INCREMENT PRIMARY KEY,
    `password` VARCHAR(15),
    `login` VARCHAR(25),
    `username` VARCHAR(20)
);

CREATE TABLE `user_infor`.`user_third`(
    `id` INT  AUTO_INCREMENT PRIMARY KEY,
    `password` VARCHAR(15),
    `login` VARCHAR(25),
    `username` VARCHAR(20)
);

CREATE TABLE `user_infor`.`user_dao`(
    `id` INT  AUTO_INCREMENT PRIMARY KEY,
    `password` VARCHAR(15),
    `login` VARCHAR(25),
    `username` VARCHAR(20),
    `status` BOOLEAN,
    `money` FLOAT
);

CREATE TABLE `user_infor`.`order`(
    `user_id` INT  AUTO_INCREMENT PRIMARY KEY,
    `order_id` INT,
    `price` INT
)
