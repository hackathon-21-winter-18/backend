DROP DATABASE IF EXISTS 21hack18;
CREATE DATABASE 21hack18;
USE 21hack18;

CREATE TABLE IF NOT EXISTS `users` (
  `username` varchar(15) NOT NULL UNIQUE,
  `hashedPass` varchar(200) NOT NULL,
  -- `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS 'palaces' (
  'id' char(36) NOT NULL,
  'name' varchar(20) NOT NULL,
  'image' VARBINARY,
  'createdBy' char(36), NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS 'palace_pins' (
  'palace' char(36), NOT NULL,
  'pin' char(36), NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;