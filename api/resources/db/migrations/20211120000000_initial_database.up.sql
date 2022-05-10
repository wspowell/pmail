-- --------------------------------------------------------
-- Host:                         snailmail-mysql.cicyow4yvieo.us-east-1.rds.amazonaws.com
-- Server version:               8.0.27 - Source distribution
-- Server OS:                    Linux
-- HeidiSQL Version:             11.2.0.6213
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Dumping structure for table snailmail_redesign.exchanges
CREATE TABLE IF NOT EXISTS `exchanges` (
  `id` int unsigned NOT NULL,
  `location_id` int unsigned NOT NULL,
  `label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_exchanges_locations` (`location_id`),
  CONSTRAINT `FK_exchanges_locations` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail_redesign.locations
CREATE TABLE IF NOT EXISTS `locations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `geo_coordinate` point NOT NULL,
  `long_1000_floor` mediumint NOT NULL,
  `lat_1000_floor` mediumint NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `address` (`address`) USING BTREE,
  KEY `lat_1000_floor` (`lat_1000_floor`),
  KEY `long_1000_floor` (`long_1000_floor`),
  SPATIAL KEY `geo_coordinate` (`geo_coordinate`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail_redesign.mail
CREATE TABLE IF NOT EXISTS `mail` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `mail_guid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `sender` int unsigned DEFAULT NULL,
  `from` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `recipient` int unsigned DEFAULT NULL,
  `to` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `starting_location` int unsigned DEFAULT NULL,
  `destination` int unsigned DEFAULT NULL,
  `address` varchar(12) COLLATE utf8mb4_unicode_ci NOT NULL,
  `contents` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `sent_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delivered_on` timestamp NULL DEFAULT NULL,
  `opened_on` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `mail_guid` (`mail_guid`) USING BTREE,
  KEY `FK_mail_users` (`sender`) USING BTREE,
  KEY `FK_mail_locations` (`destination`),
  KEY `FK_mail_users_2` (`recipient`),
  KEY `FK_mail_locations_2` (`starting_location`),
  CONSTRAINT `FK_mail_locations` FOREIGN KEY (`destination`) REFERENCES `locations` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_locations_2` FOREIGN KEY (`starting_location`) REFERENCES `locations` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_users` FOREIGN KEY (`sender`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_users_2` FOREIGN KEY (`recipient`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail_redesign.mailboxes
CREATE TABLE IF NOT EXISTS `mailboxes` (
  `location_id` int unsigned NOT NULL,
  `user_id` int unsigned NOT NULL,
  KEY `FK_mailboxes_locations` (`location_id`),
  KEY `FK_mailboxes_users` (`user_id`) USING BTREE,
  CONSTRAINT `FK_mailboxes_locations` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_mailboxes_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail_redesign.mail_tracking
CREATE TABLE IF NOT EXISTS `mail_tracking` (
  `mail_id` int unsigned NOT NULL,
  `carrier` int unsigned DEFAULT NULL,
  `exchange` int unsigned DEFAULT NULL,
  `last_carrier` int unsigned DEFAULT NULL,
  `updated_at` timestamp NOT NULL,
  KEY `FK_mail_tracking_mail` (`mail_id`),
  KEY `FK_mail_tracking_users` (`carrier`),
  KEY `FK_mail_tracking_exchanges` (`exchange`),
  KEY `FK_mail_tracking_users_2` (`last_carrier`),
  CONSTRAINT `FK_mail_tracking_exchanges` FOREIGN KEY (`exchange`) REFERENCES `exchanges` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_tracking_mail` FOREIGN KEY (`mail_id`) REFERENCES `mail` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_tracking_users` FOREIGN KEY (`carrier`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_tracking_users_2` FOREIGN KEY (`last_carrier`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail_redesign.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `guid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `public_key` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `signature` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
