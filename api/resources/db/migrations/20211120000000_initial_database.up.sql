-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               5.7.35 - MySQL Community Server (GPL)
-- Server OS:                    Linux
-- HeidiSQL Version:             11.2.0.6213
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Dumping structure for table snailmail.mail
CREATE TABLE IF NOT EXISTS `mail` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mail_guid` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from` int(10) unsigned NOT NULL,
  `to` int(10) unsigned NOT NULL,
  `contents` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `sent_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delivered_on` timestamp NULL DEFAULT NULL,
  `opened_on` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `mail_guid` (`mail_guid`),
  KEY `FK_mail_from` (`from`),
  KEY `FK_mail_to` (`to`),
  CONSTRAINT `FK_mail_from` FOREIGN KEY (`from`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_to` FOREIGN KEY (`to`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.mailboxes
CREATE TABLE IF NOT EXISTS `mailboxes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `owner` int(10) unsigned DEFAULT NULL,
  `capacity` int(10) unsigned NOT NULL,
  `latitude` decimal(10,8) NOT NULL,
  `longitude` decimal(11,8) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`) USING BTREE,
  KEY `FK_mailboxes_owner` (`owner`),
  CONSTRAINT `FK_mailboxes_owner` FOREIGN KEY (`owner`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.mailbox_mail
CREATE TABLE IF NOT EXISTS `mailbox_mail` (
  `mailbox_id` int(10) unsigned NOT NULL,
  `mail_id` int(10) unsigned NOT NULL,
  KEY `FK_mailbox_mail_mailbox_id` (`mailbox_id`),
  KEY `FK_mailbox_mail_mail_id` (`mail_id`),
  CONSTRAINT `FK_mailbox_mail_mail_id` FOREIGN KEY (`mail_id`) REFERENCES `mail` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_mailbox_mail_mailbox_id` FOREIGN KEY (`mailbox_id`) REFERENCES `mailboxes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.mail_carriers
CREATE TABLE IF NOT EXISTS `mail_carriers` (
  `user_id` int(10) unsigned NOT NULL,
  `mail_id` int(10) unsigned NOT NULL,
  UNIQUE KEY `mail_id` (`mail_id`),
  KEY `FK_mail_carriers_mail_guid` (`mail_id`) USING BTREE,
  KEY `FK_mail_carriers_user_guid` (`user_id`) USING BTREE,
  CONSTRAINT `FK_mail_carriers_mail_id` FOREIGN KEY (`mail_id`) REFERENCES `mail` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_carriers_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.mail_history
CREATE TABLE IF NOT EXISTS `mail_history` (
  `mail_id` int(10) unsigned NOT NULL,
  `mailbox_id` int(10) unsigned NOT NULL,
  `carrier` int(10) unsigned NOT NULL,
  `dropped_off_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `distance_carried` decimal(20,6) NOT NULL,
  KEY `FK_mail_history_mail_id` (`mail_id`),
  KEY `FK_mail_history_mailbox_id` (`mailbox_id`),
  KEY `FK_mail_history_carrier` (`carrier`),
  CONSTRAINT `FK_mail_history_carrier` FOREIGN KEY (`carrier`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_history_mail_id` FOREIGN KEY (`mail_id`) REFERENCES `mail` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_mail_history_mailbox_id` FOREIGN KEY (`mailbox_id`) REFERENCES `mailboxes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.penpals
CREATE TABLE IF NOT EXISTS `penpals` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `sender` int(10) unsigned NOT NULL,
  `receiver` int(10) unsigned NOT NULL,
  `pending` tinyint(1) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `dRelationships_status` (`pending`) USING BTREE,
  KEY `FK_penpals_receiver` (`receiver`),
  KEY `FK_penpals_sender` (`sender`),
  CONSTRAINT `FK_penpals_receiver` FOREIGN KEY (`receiver`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_penpals_sender` FOREIGN KEY (`sender`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.schema_migrations
CREATE TABLE IF NOT EXISTS `schema_migrations` (
  `version` bigint(20) NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_guid` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `pineapple_on_pizza` tinyint(1) unsigned NOT NULL,
  `mailbag_capacity` int(10) unsigned NOT NULL,
  `created_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `secret` char(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `salt` char(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `username` (`username`) USING BTREE,
  UNIQUE KEY `user_guid` (`user_guid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.user_blocks
CREATE TABLE IF NOT EXISTS `user_blocks` (
  `blocker` int(10) unsigned NOT NULL,
  `blockee` int(10) unsigned NOT NULL,
  UNIQUE KEY `blocker_blockee` (`blocker`,`blockee`),
  KEY `blockee` (`blockee`) USING BTREE,
  KEY `blocker` (`blocker`),
  CONSTRAINT `FK_user_blocks_blockee` FOREIGN KEY (`blockee`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_user_blocks_blocker` FOREIGN KEY (`blocker`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.user_inbox
CREATE TABLE IF NOT EXISTS `user_inbox` (
  `mail_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`mail_id`),
  KEY `FK_user_inbox_mail_id` (`mail_id`),
  CONSTRAINT `FK_user_inbox_mail_id` FOREIGN KEY (`mail_id`) REFERENCES `mail` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.user_mailbox_labels
CREATE TABLE IF NOT EXISTS `user_mailbox_labels` (
  `user_id` int(10) unsigned DEFAULT NULL,
  `mailbox_id` int(10) unsigned DEFAULT NULL,
  `label` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  UNIQUE KEY `user_id_mailbox_id` (`user_id`,`mailbox_id`) USING BTREE,
  KEY `user_id` (`user_id`),
  KEY `mailbox_id` (`mailbox_id`),
  CONSTRAINT `FK_user_mailbox_labels_mailbox_id` FOREIGN KEY (`mailbox_id`) REFERENCES `mailboxes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_user_mailbox_labels_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.user_outbox
CREATE TABLE IF NOT EXISTS `user_outbox` (
  `mail_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`mail_id`),
  KEY `FK_user_outbox_mail_id` (`mail_id`),
  CONSTRAINT `FK_user_outbox_mail_id` FOREIGN KEY (`mail_id`) REFERENCES `mail` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

-- Dumping structure for table snailmail.user_penpals
CREATE TABLE IF NOT EXISTS `user_penpals` (
  `user_id` int(11) unsigned NOT NULL,
  `penpal_id` int(11) unsigned NOT NULL,
  UNIQUE KEY `user_id_penpal_id` (`user_id`,`penpal_id`),
  KEY `penpal_id` (`penpal_id`) USING BTREE,
  KEY `user_id` (`user_id`) USING BTREE,
  CONSTRAINT `FK_user_penpals_penpal_id` FOREIGN KEY (`penpal_id`) REFERENCES `penpals` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_user_penpals_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Data exporting was unselected.

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
