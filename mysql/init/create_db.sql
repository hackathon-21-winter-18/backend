DROP DATABASE IF EXISTS 21hack18;
CREATE DATABASE 21hack18;
USE 21hack18;

CREATE TABLE IF NOT EXISTS `users` (
  `id` char(36) NOT NULL UNIQUE,
  `googleID` char(21) NOT NUll UNIQUE,
  `name` varchar(20) NOT NULL UNIQUE,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `palaces` (
  `id` char(36) NOT NULL UNIQUE,
  `originalID` char(36) NOT NUll,
  `name` varchar(20) NOT NULL,
  `createdBy` char(36) NOT NULL,
  `image` varchar(40) NOT NULL,
  `heldBy` char(36) NOT NULL,
  `number_of_embededPins` int DEFAULT 0,
  `share` boolean DEFAULT False,
  `savedCount` int DEFAULT 0,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NUll, 
  `shared_at` datetime NULL, 
  `firstshared` boolean DEFAULT False,
  `firstshared_at` datetime NULL,
  `group1` varchar(20) DEFAULT "",
  `group2` varchar(20) DEFAULT "",
  `group3` varchar(20) DEFAULT "",
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
  `situation` varchar(15) NULL,
  `palaceID` char(36) NOT NULL,
  `groupNumber` int DEFAULT 0,
  FOREIGN KEY (`palaceID`) REFERENCES palaces(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `templates` (
  `id` char(36) NOT NULL UNIQUE,
  `originalID` char(36) NOT NUll,
  `name` varchar(20) NOT NULL,
  `createdBy` char(36) NOT NULL,
  `image` varchar(40) NOT NULL,
  `heldBy` char(36) NULL,
  `number_of_pins` int DEFAULT 0,
  `share` boolean DEFAULT False,
  `savedCount` int DEFAULT 0,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NULL,
  `shared_at` datetime NULL,
  `firstshared` boolean DEFAULT False,
  `firstshared_at` datetime NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`createdBy`) REFERENCES users(`id`),
  FOREIGN KEY (`heldBy`) REFERENCES users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `pins` (
  `number` int NOT NULL,
  `x` decimal(10, 2) NOT NULL, 
  `y` decimal(10, 2) NOT NULL,
  `templateID` char(36) NOT NULL,
  `groupNumber` int DEFAULT 0,
  FOREIGN KEY (`templateID`) REFERENCES templates(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `palace_user` (
  `palaceID` char(36) NOT NUll,
  `userID` char(36) NOT NUll,
  FOREIGN KEY (`palaceID`) REFERENCES palaces(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`userID`) REFERENCES users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `template_user` (
  `templateID` char(36) NOT NUll,
  `userID` char(36) NOT NUll,
  FOREIGN KEY (`templateID`) REFERENCES templates(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`userID`) REFERENCES users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `notices` (
  `id` char(36) NOT NULL,
  `userID` char(36) NOT NULL,
  `content` varchar(400) NOT NULL,
  `checked` boolean DEFAULT False,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `palaceID` char(36) NULL,
  `templateID` char(36) NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`userID`) REFERENCES users(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

