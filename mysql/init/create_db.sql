DROP DATABASE IF EXISTS 21hack18;
CREATE DATABASE 21hack18;
USE 21hack18;

CREATE TABLE IF NOT EXISTS `users` (
  `id` char(36) NOT NULL UNIQUE,
  `name` varchar(15) NOT NULL UNIQUE,
  `hashedPass` varchar(200) NOT NULL,
  -- `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `palaces` (
  `id` char(36) NOT NULL UNIQUE,
  `name` varchar(20) NOT NULL,
  `createdBy` char(36) NOT NULL,
  `image` varchar(40) NOT NULL,
  `heldBy` char(36) NOT NULL,
  `share` boolean DEFAULT False,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NUll, 
  `shared_at` datetime NULL, 
  `firstshared` boolean DEFAULT False,
  `firstshared_at` datetime NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`createdBy`) REFERENCES users(`id`),
  FOREIGN KEY (`heldBy`) REFERENCES users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `embededpins` (
  `number` int NOT NULL,
  `x` decimal(10, 2) NOT NULL, 
  `y` decimal(10, 2) NOT NULL,
  `word` varchar(15) NULL,
  `place` varchar(15) NULL,
  `do` varchar(15) NULL,
  `palaceID` char(36) NOT NULL,
  FOREIGN KEY (`palaceID`) REFERENCES palaces(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `templates` (
  `id` char(36) NOT NULL UNIQUE,
  `name` varchar(20) NOT NULL,
  `createdBy` char(36) NOT NULL,
  `image` varchar(40) NOT NULL,
  `heldBy` char(36) NULL,
  `share` boolean DEFAULT False,  
  `created_at` datetime NOT NULL,
  `updated_at` datetime NULL,
  `shared_at` datetime NULL,
  `firstshared` boolean DEFAULT False,
  `firstshared_at` datetime NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`createdBy`) REFERENCES users(`id`),
  FOREIGN KEY (`heldBy`) REFERENCES users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `templatepins` (
  `number` int NOT NULL,
  `x` decimal(10, 2) NOT NULL, 
  `y` decimal(10, 2) NOT NULL,
  `templateID` char(36) NOT NULL,
  FOREIGN KEY (`templateID`) REFERENCES templates(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;