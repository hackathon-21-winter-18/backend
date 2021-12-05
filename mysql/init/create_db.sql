DROP DATABASE IF EXISTS 21hack18;
CREATE DATABASE 21hack18;
USE 21hack18;

CREATE TABLE IF NOT EXISTS `palaces` (
  `id` char(36) NOT NULL,
  `name` varchar(20) NOT NULL,
  `image` VARBINARY(100),
  `createdBy` char(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `pin` (
  `id` char(36) NOT NULL,
  `x` int NOT NULL,
  `y` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `palace_pins` (
  `palace` char(36) NOT NULL,
  `pin` char(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `word` (
  `id` char(36) NOT NULL,
  `pin` char(36) NOT NULL,
  `word` varchar(36) NOT NULL,
  `memo` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `palace_words` (
  `palace` char(36) NOT NULL,
  `word` char(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `template` (
  `id` char(36) NOT NULL,
  `name` varchar(20) NOT NULL,
  `image` VARBINARY(100),
  `createdBy` char(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `template_pins` (
  `template` char(36) NOT NULL,
  `pin` char(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user` (
 `id` char(36) NOT NULL UNIQUE,
  `username` varchar(15) NOT NULL,
  `hashedPass` varchar(200) NOT NULL,
  PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_palaces` (
  `user`  char(36) NOT NULL,
  `palace`char(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;