CREATE TABLE IF NOT EXISTS `sport_rank` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `data` text COLLATE utf8mb4_general_ci NOT NULL,
  `date` date NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT (now()) ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`) USING BTREE,
  KEY `idx_type_date` (`type`,`date`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
