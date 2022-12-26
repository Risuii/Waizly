CREATE TABLE `waizly`.`account` (
  `ID` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(255) NULL,
  `password` VARCHAR(255) NULL,
  `email` VARCHAR(255) NULL,
  `created_at` DATETIME NULL DEFAULT (now()),
  `update_at` DATETIME NULL DEFAULT (now()),
  PRIMARY KEY (`ID`)
);