CREATE TABLE IF NOT EXISTS `news` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '標題',
  `description` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '描述',
  `cover` varchar(155) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '封面',
  `cover_source` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '原始封面連結',
  `link` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '新聞連結',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '內文',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '狀態，上架:1、下架:0',
  `source` varchar(30) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '來源',
  `pub_date` timestamp NOT NULL COMMENT '發佈時間',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '建立時間',
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新時間',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `unique_title_source` (`title`, `source`) USING BTREE,
  KEY `idx_title_source` (`title`,`source`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
