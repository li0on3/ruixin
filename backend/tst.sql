-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: ruixin_platform
-- ------------------------------------------------------
-- Server version	8.0.42-0ubuntu0.24.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `admin_operation_logs`
--

DROP TABLE IF EXISTS `admin_operation_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin_operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `admin_id` bigint unsigned DEFAULT NULL,
  `admin_name` longtext,
  `operation` longtext,
  `module` longtext,
  `details` text,
  `ip_address` longtext,
  `user_agent` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin_operation_logs`
--

LOCK TABLES `admin_operation_logs` WRITE;
/*!40000 ALTER TABLE `admin_operation_logs` DISABLE KEYS */;
INSERT INTO `admin_operation_logs` VALUES (1,'2025-05-24 09:08:01.955',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(2,'2025-05-25 12:53:30.186',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(3,'2025-05-25 13:00:08.644',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(4,'2025-05-25 13:05:18.284',1,'','reset_api_key','distributor','Reset API key for distributor: 示例分销商','::1','PostmanRuntime/7.44.0'),(5,'2025-05-25 13:21:28.397',1,'','create_card','card','Created card: VTEX97M4','::1','PostmanRuntime/7.44.0'),(6,'2025-05-25 15:45:19.257',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(7,'2025-05-25 15:49:36.433',1,'','logout','auth','Admin logout','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(8,'2025-05-25 15:49:51.518',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(9,'2025-05-25 16:11:12.779',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(10,'2025-05-25 16:17:06.106',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(11,'2025-05-25 16:17:30.350',1,'','create_distributor','distributor','Created distributor: 新分销商','::1','PostmanRuntime/7.44.0'),(12,'2025-05-25 18:09:19.522',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(13,'2025-05-25 18:11:15.540',1,'','create_distributor','distributor','Created distributor: 测试分销商','::1','PostmanRuntime/7.44.0'),(14,'2025-05-25 18:21:36.414',1,'','create_card','card','Created card: 8HVXKL76','::1','PostmanRuntime/7.44.0'),(15,'2025-05-25 18:32:22.064',1,'','create_card','card','Created card: YQC66PDG','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(16,'2025-05-25 18:37:31.256',1,'','create_card','card','Created card: KK8NY7YK','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(17,'2025-05-25 20:33:13.822',1,'','logout','auth','Admin logout','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(18,'2025-05-25 20:33:47.783',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(19,'2025-05-27 07:40:47.475',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(20,'2025-05-27 07:40:47.850',1,'','logout','auth','Admin logout','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(21,'2025-05-27 07:40:56.741',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(22,'2025-05-27 22:43:42.545',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(23,'2025-05-28 20:22:19.509',1,'','batch_validate_cards','card','Validated 4 cards','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(24,'2025-05-28 20:23:32.533',1,'','batch_validate_cards','card','Validated 0 cards','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(25,'2025-05-28 20:23:49.117',1,'','delete_card','card','Deleted card: 8HVXKL76','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(26,'2025-05-28 20:42:33.113',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(27,'2025-05-28 21:03:00.848',1,'','delete_card','card','Deleted card ID: 25','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(28,'2025-05-28 22:49:51.470',1,'','delete_card','card','Deleted card ID: 6','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(29,'2025-05-28 22:49:54.841',1,'','delete_card','card','Deleted card ID: 5','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(30,'2025-05-28 22:56:01.119',1,'','create_card','card','Created card: YQC66PDG','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(31,'2025-05-28 23:21:17.868',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(32,'2025-05-29 07:30:19.024',1,'','create_card','card','Created card: KSGHGK4C','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(33,'2025-05-29 19:45:21.771',1,'admin','login','auth','Admin login successful','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(34,'2025-05-29 19:46:06.532',1,'','delete_card','card','Deleted card ID: 28','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(35,'2025-05-29 19:46:22.435',1,'','create_card','card','Created card: KSGHGK4C','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(36,'2025-05-29 19:47:27.128',1,'','batch_import_cards','card','Batch imported 10 cards','127.0.0.1','Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36'),(37,'2025-05-29 20:01:12.772',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(38,'2025-05-29 20:34:50.692',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(39,'2025-05-29 21:31:58.869',1,'admin','login','auth','Admin login successful','::1','PostmanRuntime/7.44.0'),(40,'2025-05-29 23:32:10.401',1,'admin','login','auth','Admin login successful','127.0.0.1','Go-http-client/1.1'),(41,'2025-05-29 23:36:28.207',1,'admin','login','auth','Admin login successful','127.0.0.1','Go-http-client/1.1');
/*!40000 ALTER TABLE `admin_operation_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `admins`
--

DROP TABLE IF EXISTS `admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admins` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(100) NOT NULL,
  `password` varchar(255) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `real_name` varchar(50) DEFAULT NULL,
  `role` varchar(50) DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `last_login_at` datetime(3) DEFAULT NULL,
  `last_login_ip` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  UNIQUE KEY `idx_admins_email` (`email`),
  KEY `idx_admins_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admins`
--

LOCK TABLES `admins` WRITE;
/*!40000 ALTER TABLE `admins` DISABLE KEYS */;
INSERT INTO `admins` VALUES (1,'2025-05-24 08:57:16.735','2025-05-29 23:36:28.199',NULL,'admin','$2a$10$Tz.GcgLXFpke/WzxmskBeOyMbsoA6n6t3ZdhWUlTDL7/n2Bba.Hlu','admin@ruixin.com','','系统管理员','super_admin',1,'2025-05-29 23:36:28.000','127.0.0.1');
/*!40000 ALTER TABLE `admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `card_batches`
--

DROP TABLE IF EXISTS `card_batches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `card_batches` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `batch_no` varchar(50) NOT NULL,
  `price_id` bigint unsigned DEFAULT NULL,
  `cost_price` decimal(10,2) DEFAULT NULL,
  `sell_price` decimal(10,2) DEFAULT NULL,
  `total_count` bigint DEFAULT NULL,
  `used_count` bigint DEFAULT NULL,
  `imported_at` datetime(3) DEFAULT NULL,
  `imported_by` bigint unsigned DEFAULT NULL,
  `description` text,
  `luckin_product_id` bigint DEFAULT '6',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_card_batches_batch_no` (`batch_no`),
  KEY `idx_card_batches_deleted_at` (`deleted_at`),
  KEY `fk_card_batches_admin` (`imported_by`),
  KEY `idx_card_batches_price_id` (`price_id`),
  CONSTRAINT `fk_card_batches_admin` FOREIGN KEY (`imported_by`) REFERENCES `admins` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `card_batches`
--

LOCK TABLES `card_batches` WRITE;
/*!40000 ALTER TABLE `card_batches` DISABLE KEYS */;
INSERT INTO `card_batches` VALUES (1,'2025-05-27 19:44:34.909','2025-05-27 19:44:34.909',NULL,'BATCH-2025-001',0,7.50,8.50,5,0,'2025-05-27 19:44:34.909',1,'示例批次',6),(2,'2025-05-28 20:59:41.612','2025-05-28 20:59:41.612',NULL,'TEST-20250528205941',1,5.00,9.10,5,0,'2025-05-28 20:59:41.612',1,'测试批次',6),(3,'2025-05-29 19:47:27.003','2025-05-29 19:47:27.003',NULL,'BATCH-20250529-9247',1,5.00,10.00,10,0,'2025-05-29 19:47:27.003',1,'',6);
/*!40000 ALTER TABLE `card_batches` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `card_product_bindings`
--

DROP TABLE IF EXISTS `card_product_bindings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `card_product_bindings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `card_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `priority` bigint DEFAULT '0',
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `idx_card_product_bindings_deleted_at` (`deleted_at`),
  KEY `idx_card_product_bindings_card_id` (`card_id`),
  KEY `idx_card_product_bindings_product_id` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `card_product_bindings`
--

LOCK TABLES `card_product_bindings` WRITE;
/*!40000 ALTER TABLE `card_product_bindings` DISABLE KEYS */;
INSERT INTO `card_product_bindings` VALUES (1,'2025-05-28 21:19:39.204','2025-05-28 21:19:39.204',NULL,5,1,0,1),(2,'2025-05-28 21:19:39.212','2025-05-28 21:19:39.212',NULL,5,2,0,1),(3,'2025-05-28 21:19:39.221','2025-05-28 21:19:39.221',NULL,5,3,0,1),(4,'2025-05-28 21:19:39.230','2025-05-28 21:19:39.230',NULL,5,4,0,1),(5,'2025-05-28 22:31:41.501','2025-05-28 22:31:41.501',NULL,6,1,0,1),(6,'2025-05-28 22:31:41.681','2025-05-28 22:31:41.681',NULL,6,2,0,1),(7,'2025-05-28 22:31:41.939','2025-05-28 22:31:41.939',NULL,6,3,0,1),(8,'2025-05-28 22:31:42.086','2025-05-28 22:31:42.086',NULL,6,4,0,1),(9,'2025-05-28 22:56:02.711','2025-05-28 22:56:02.711',NULL,27,1,0,1),(10,'2025-05-28 22:56:02.911','2025-05-28 22:56:02.911',NULL,27,2,0,1),(11,'2025-05-28 22:56:03.180','2025-05-28 22:56:03.180',NULL,27,3,0,1),(12,'2025-05-28 22:56:03.347','2025-05-28 22:56:03.347',NULL,27,4,0,1),(13,'2025-05-29 07:30:20.563','2025-05-29 07:30:20.563',NULL,28,1,0,1),(14,'2025-05-29 07:30:20.738','2025-05-29 07:30:20.738',NULL,28,2,0,1),(15,'2025-05-29 07:30:21.017','2025-05-29 07:30:21.017',NULL,28,3,0,1),(16,'2025-05-29 07:30:21.172','2025-05-29 07:30:21.172',NULL,28,4,0,1),(17,'2025-05-29 19:46:24.151','2025-05-29 19:46:24.151',NULL,29,1,0,1),(18,'2025-05-29 19:46:24.371','2025-05-29 19:46:24.371',NULL,29,2,0,1),(19,'2025-05-29 19:46:24.672','2025-05-29 19:46:24.672',NULL,29,3,0,1),(20,'2025-05-29 19:46:24.837','2025-05-29 19:46:24.837',NULL,29,4,0,1),(21,'2025-05-29 19:47:28.697','2025-05-29 19:47:28.697',NULL,31,1,0,1),(22,'2025-05-29 19:47:28.905','2025-05-29 19:47:28.905',NULL,31,2,0,1),(23,'2025-05-29 19:47:29.200','2025-05-29 19:47:29.200',NULL,31,3,0,1),(24,'2025-05-29 19:47:29.362','2025-05-29 19:47:29.362',NULL,31,4,0,1),(25,'2025-05-29 19:47:35.331','2025-05-29 19:47:35.331',NULL,35,1,0,1),(26,'2025-05-29 19:47:35.513','2025-05-29 19:47:35.513',NULL,35,2,0,1),(27,'2025-05-29 19:47:35.791','2025-05-29 19:47:35.791',NULL,35,3,0,1),(28,'2025-05-29 19:47:35.953','2025-05-29 19:47:35.953',NULL,35,4,0,1);
/*!40000 ALTER TABLE `card_product_bindings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `card_usage_logs`
--

DROP TABLE IF EXISTS `card_usage_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `card_usage_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `card_id` bigint unsigned DEFAULT NULL,
  `distributor_id` bigint unsigned DEFAULT NULL,
  `order_id` longtext,
  `success` tinyint(1) DEFAULT NULL,
  `error_message` text,
  `card_code` varchar(100) DEFAULT NULL,
  `order_no` varchar(50) DEFAULT NULL,
  `fail_reason` varchar(50) DEFAULT NULL,
  `request_ip` varchar(45) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_card_usage_logs_card` (`card_id`),
  KEY `fk_card_usage_logs_distributor` (`distributor_id`),
  CONSTRAINT `fk_card_usage_logs_card` FOREIGN KEY (`card_id`) REFERENCES `cards` (`id`),
  CONSTRAINT `fk_card_usage_logs_distributor` FOREIGN KEY (`distributor_id`) REFERENCES `distributors` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `card_usage_logs`
--

LOCK TABLES `card_usage_logs` WRITE;
/*!40000 ALTER TABLE `card_usage_logs` DISABLE KEYS */;
INSERT INTO `card_usage_logs` VALUES (1,'2025-05-28 23:23:44.424',27,3,NULL,1,'','YQC66PDG','DD174844582354738625164','','',''),(2,'2025-05-29 21:56:04.017',29,3,NULL,1,'','KSGHGK4C','DD1748526961752133ad737','','','');
/*!40000 ALTER TABLE `card_usage_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cards`
--

DROP TABLE IF EXISTS `cards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cards` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `card_code` varchar(100) NOT NULL,
  `luckin_product_id` bigint DEFAULT '6',
  `status` bigint DEFAULT '0',
  `product_id` bigint DEFAULT NULL,
  `daily_limit` bigint DEFAULT NULL,
  `total_limit` bigint DEFAULT NULL,
  `used_count` bigint DEFAULT NULL,
  `expired_at` datetime(3) DEFAULT NULL,
  `description` text,
  `batch_id` bigint unsigned DEFAULT NULL,
  `price_id` bigint DEFAULT NULL,
  `cost_price` decimal(10,2) DEFAULT NULL,
  `sell_price` decimal(10,2) DEFAULT NULL,
  `used_at` datetime(3) DEFAULT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  `reserved_at` datetime(3) DEFAULT NULL,
  `sync_status` varchar(20) DEFAULT 'pending',
  `synced_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cards_deleted_at` (`deleted_at`),
  KEY `fk_cards_order` (`order_id`),
  KEY `fk_card_batches_cards` (`batch_id`),
  KEY `idx_luckin_product_id` (`luckin_product_id`),
  CONSTRAINT `fk_card_batches_cards` FOREIGN KEY (`batch_id`) REFERENCES `card_batches` (`id`),
  CONSTRAINT `fk_cards_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cards`
--

LOCK TABLES `cards` WRITE;
/*!40000 ALTER TABLE `cards` DISABLE KEYS */;
INSERT INTO `cards` VALUES (1,'2025-05-24 08:57:16.746','2025-05-30 01:35:14.497',NULL,'DEMO123456',6,1,6,10,100,0,'2026-05-24 08:57:16.745','示例优惠卡',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'synced','2025-05-30 01:35:14.298'),(2,'2025-05-25 13:21:28.388','2025-05-30 01:35:14.639',NULL,'VTEX97M4',6,1,6,10,100,0,'2030-01-01 07:59:59.000','测试卡片',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'synced','2025-05-30 01:35:14.298'),(6,'2025-05-25 18:37:31.244','2025-05-28 22:46:31.776','2025-05-28 22:49:51.459','KK8NY7YK',6,0,6,0,100,0,'2026-05-25 18:37:10.163','',NULL,NULL,NULL,NULL,NULL,NULL,NULL,'synced','2025-05-28 22:46:31.034'),(8,'2025-05-27 19:44:34.922','2025-05-30 01:35:14.762',NULL,'DEMO123457',6,1,NULL,NULL,NULL,NULL,'2026-05-27 19:44:34.918','',1,0,7.50,8.50,'2025-05-28 20:22:18.840',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(9,'2025-05-27 19:44:34.927','2025-05-30 01:35:14.887',NULL,'DEMO123458',6,1,NULL,NULL,NULL,NULL,'2026-05-27 19:44:34.918','',1,0,7.50,8.50,'2025-05-28 20:22:18.840',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(10,'2025-05-27 19:44:34.938','2025-05-30 01:35:15.011',NULL,'DEMO123459',6,1,NULL,NULL,NULL,NULL,'2026-05-27 19:44:34.918','',1,0,7.50,8.50,'2025-05-28 20:22:18.840',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(11,'2025-05-27 19:44:34.947','2025-05-30 01:35:15.133',NULL,'DEMO123460',6,1,NULL,NULL,NULL,NULL,'2026-05-27 19:44:34.918','',1,0,7.50,8.50,'2025-05-28 20:22:18.840',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(21,'2025-05-28 20:59:41.624','2025-05-30 01:35:15.255',NULL,'TEST001',6,1,NULL,NULL,NULL,NULL,'2025-11-28 20:59:41.623','测试卡片001',2,1,5.00,9.10,'2025-05-28 21:01:32.274',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(22,'2025-05-28 20:59:41.632','2025-05-30 01:35:15.383',NULL,'TEST002',6,1,NULL,NULL,NULL,NULL,'2025-11-28 20:59:41.623','测试卡片002',2,1,5.00,9.10,'2025-05-28 21:01:32.274',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(23,'2025-05-28 20:59:41.640','2025-05-30 01:35:15.502',NULL,'TEST003',6,1,NULL,NULL,NULL,NULL,'2025-11-28 20:59:41.623','测试卡片003',2,1,5.00,9.10,'2025-05-28 21:01:32.274',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(24,'2025-05-28 20:59:41.649','2025-05-30 01:35:15.630',NULL,'TEST004',6,1,NULL,NULL,NULL,NULL,'2025-11-28 20:59:41.623','测试卡片004 - 已使用',2,1,5.00,9.10,'2025-05-27 20:59:41.623',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(25,'2025-05-28 20:59:41.658','2025-05-28 20:59:41.658','2025-05-28 21:03:00.838','TEST005',6,2,NULL,NULL,NULL,NULL,'2025-11-28 20:59:41.623','测试卡片005 - 预占中',2,1,5.00,9.10,NULL,NULL,'2025-05-28 19:59:41.623','pending',NULL),(27,'2025-05-28 22:56:01.110','2025-05-30 01:35:16.783',NULL,'YQC66PDG',6,1,NULL,NULL,NULL,NULL,'2026-05-28 22:49:55.523','',NULL,6,5.00,10.00,'2025-05-29 01:09:59.895',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(29,'2025-05-29 19:46:22.423','2025-05-30 01:35:16.922',NULL,'KSGHGK4C',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:07.610','',NULL,1,5.00,10.00,'2025-05-30 01:35:14.298',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(30,'2025-05-29 19:47:27.014','2025-05-30 01:35:17.061',NULL,'4H8SB644',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 19:47:32.845',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(31,'2025-05-29 19:47:27.028','2025-05-30 01:35:17.208',NULL,'8HVXKL76',6,0,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,NULL,NULL,NULL,'synced','2025-05-30 01:35:14.298'),(32,'2025-05-29 19:47:27.042','2025-05-30 01:35:17.346',NULL,'KSGHGK4C',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 22:00:59.160',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(33,'2025-05-29 19:47:27.056','2025-05-30 01:35:17.488',NULL,'YQC66PDG',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 19:47:32.845',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(34,'2025-05-29 19:47:27.069','2025-05-30 01:35:17.625',NULL,'VTEX97M4',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 19:47:32.845',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(35,'2025-05-29 19:47:27.079','2025-05-30 01:35:17.759',NULL,'KK8NY7YK',6,0,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,NULL,NULL,NULL,'synced','2025-05-30 01:35:14.298'),(36,'2025-05-29 19:47:27.089','2025-05-30 01:35:17.899',NULL,'VDCU4V2Z',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 19:47:32.845',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(37,'2025-05-29 19:47:27.099','2025-05-30 01:35:18.034',NULL,'HWFG3RLL',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 19:47:32.845',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(38,'2025-05-29 19:47:27.111','2025-05-30 01:35:19.174',NULL,'N58DF6SV',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 19:47:32.845',NULL,NULL,'synced','2025-05-30 01:35:14.298'),(39,'2025-05-29 19:47:27.118','2025-05-30 01:35:19.312',NULL,'YS5QNFH9',6,1,NULL,NULL,NULL,NULL,'2026-05-29 19:46:55.024','',3,1,5.00,10.00,'2025-05-29 19:47:32.845',NULL,NULL,'synced','2025-05-30 01:35:14.298');
/*!40000 ALTER TABLE `cards` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `category_bindings`
--

DROP TABLE IF EXISTS `category_bindings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `category_bindings` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned NOT NULL,
  `target_type` varchar(20) NOT NULL,
  `target_id` varchar(50) NOT NULL,
  `priority` bigint DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_bindings_category_id` (`category_id`),
  KEY `fk_category_bindings_creator` (`created_by`),
  CONSTRAINT `fk_category_bindings_category` FOREIGN KEY (`category_id`) REFERENCES `cards` (`id`),
  CONSTRAINT `fk_category_bindings_creator` FOREIGN KEY (`created_by`) REFERENCES `admins` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `category_bindings`
--

LOCK TABLES `category_bindings` WRITE;
/*!40000 ALTER TABLE `category_bindings` DISABLE KEYS */;
/*!40000 ALTER TABLE `category_bindings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cities`
--

DROP TABLE IF EXISTS `cities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `city_id` bigint NOT NULL,
  `city_name` varchar(50) NOT NULL,
  `pinyin` varchar(10) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_cities_city_id` (`city_id`),
  KEY `idx_cities_city_name` (`city_name`)
) ENGINE=InnoDB AUTO_INCREMENT=363 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cities`
--

LOCK TABLES `cities` WRITE;
/*!40000 ALTER TABLE `cities` DISABLE KEYS */;
INSERT INTO `cities` VALUES (1,1,'北京','B','2025-05-25 18:21:49.715','2025-05-25 18:21:49.715'),(2,2,'天津','T','2025-05-25 18:21:49.715','2025-05-25 18:21:49.715'),(3,3,'石家庄','S','2025-05-25 18:21:49.716','2025-05-25 18:21:49.716'),(4,4,'唐山','T','2025-05-25 18:21:49.716','2025-05-25 18:21:49.716'),(5,5,'秦皇岛','Q','2025-05-25 18:21:49.716','2025-05-25 18:21:49.716'),(6,6,'邯郸','H','2025-05-25 18:21:49.716','2025-05-25 18:21:49.716'),(7,7,'保定','B','2025-05-25 18:21:49.717','2025-05-25 18:21:49.717'),(8,8,'廊坊','L','2025-05-25 18:21:49.717','2025-05-25 18:21:49.717'),(9,9,'张家口','Z','2025-05-25 18:21:49.717','2025-05-25 18:21:49.717'),(10,10,'沧州','C','2025-05-25 18:21:49.718','2025-05-25 18:21:49.718'),(11,11,'邢台','X','2025-05-25 18:21:49.718','2025-05-25 18:21:49.718'),(12,12,'承德','C','2025-05-25 18:21:49.718','2025-05-25 18:21:49.718'),(13,13,'衡水','H','2025-05-25 18:21:49.718','2025-05-25 18:21:49.718'),(14,14,'太原','T','2025-05-25 18:21:49.719','2025-05-25 18:21:49.719'),(15,15,'临汾','L','2025-05-25 18:21:49.719','2025-05-25 18:21:49.719'),(16,16,'晋城','J','2025-05-25 18:21:49.719','2025-05-25 18:21:49.719'),(17,17,'运城','Y','2025-05-25 18:21:49.719','2025-05-25 18:21:49.719'),(18,18,'大同','D','2025-05-25 18:21:49.720','2025-05-25 18:21:49.720'),(19,19,'阳泉','Y','2025-05-25 18:21:49.720','2025-05-25 18:21:49.720'),(20,20,'长治','C','2025-05-25 18:21:49.720','2025-05-25 18:21:49.720'),(21,21,'晋中','J','2025-05-25 18:21:49.720','2025-05-25 18:21:49.720'),(22,22,'吕梁','L','2025-05-25 18:21:49.721','2025-05-25 18:21:49.721'),(23,23,'朔州','S','2025-05-25 18:21:49.721','2025-05-25 18:21:49.721'),(24,24,'忻州','X','2025-05-25 18:21:49.721','2025-05-25 18:21:49.721'),(25,25,'呼和浩特','H','2025-05-25 18:21:49.721','2025-05-25 18:21:49.721'),(26,26,'包头','B','2025-05-25 18:21:49.722','2025-05-25 18:21:49.722'),(27,27,'赤峰','C','2025-05-25 18:21:49.722','2025-05-25 18:21:49.722'),(28,28,'鄂尔多斯','E','2025-05-25 18:21:49.722','2025-05-25 18:21:49.722'),(29,29,'呼伦贝尔','H','2025-05-25 18:21:49.723','2025-05-25 18:21:49.723'),(30,30,'锡林郭勒','X','2025-05-25 18:21:49.723','2025-05-25 18:21:49.723'),(31,31,'巴彦淖尔','B','2025-05-25 18:21:49.723','2025-05-25 18:21:49.723'),(32,32,'通辽','T','2025-05-25 18:21:49.724','2025-05-25 18:21:49.724'),(33,33,'乌兰察布','W','2025-05-25 18:21:49.724','2025-05-25 18:21:49.724'),(34,34,'乌海','W','2025-05-25 18:21:49.724','2025-05-25 18:21:49.724'),(35,35,'阿拉善盟','A','2025-05-25 18:21:49.724','2025-05-25 18:21:49.724'),(36,36,'兴安盟','X','2025-05-25 18:21:49.725','2025-05-25 18:21:49.725'),(37,37,'乌兰浩特','W','2025-05-25 18:21:49.725','2025-05-25 18:21:49.725'),(38,38,'沈阳','S','2025-05-25 18:21:49.725','2025-05-25 18:21:49.725'),(39,39,'大连','D','2025-05-25 18:21:49.726','2025-05-25 18:21:49.726'),(40,40,'鞍山','A','2025-05-25 18:21:49.726','2025-05-25 18:21:49.726'),(41,41,'抚顺','F','2025-05-25 18:21:49.726','2025-05-25 18:21:49.726'),(42,42,'锦州','J','2025-05-25 18:21:49.726','2025-05-25 18:21:49.726'),(43,43,'葫芦岛','H','2025-05-25 18:21:49.727','2025-05-25 18:21:49.727'),(44,44,'辽阳','L','2025-05-25 18:21:49.727','2025-05-25 18:21:49.727'),(45,45,'盘锦','P','2025-05-25 18:21:49.727','2025-05-25 18:21:49.727'),(46,46,'营口','Y','2025-05-25 18:21:49.727','2025-05-25 18:21:49.727'),(47,47,'丹东','D','2025-05-25 18:21:49.728','2025-05-25 18:21:49.728'),(48,48,'阜新','F','2025-05-25 18:21:49.728','2025-05-25 18:21:49.728'),(49,49,'铁岭','T','2025-05-25 18:21:49.728','2025-05-25 18:21:49.728'),(50,50,'朝阳','Z','2025-05-25 18:21:49.728','2025-05-25 18:21:49.728'),(51,51,'本溪','B','2025-05-25 18:21:49.729','2025-05-25 18:21:49.729'),(52,52,'长春','C','2025-05-25 18:21:49.729','2025-05-25 18:21:49.729'),(53,53,'吉林','J','2025-05-25 18:21:49.729','2025-05-25 18:21:49.729'),(54,54,'白山','B','2025-05-25 18:21:49.729','2025-05-25 18:21:49.729'),(55,55,'通化','T','2025-05-25 18:21:49.729','2025-05-25 18:21:49.729'),(56,56,'延边','Y','2025-05-25 18:21:49.729','2025-05-25 18:21:49.729'),(57,57,'四平','S','2025-05-25 18:21:49.730','2025-05-25 18:21:49.730'),(58,58,'松原','S','2025-05-25 18:21:49.730','2025-05-25 18:21:49.730'),(59,59,'白城','B','2025-05-25 18:21:49.730','2025-05-25 18:21:49.730'),(60,60,'辽源','L','2025-05-25 18:21:49.730','2025-05-25 18:21:49.730'),(61,61,'哈尔滨','H','2025-05-25 18:21:49.731','2025-05-25 18:21:49.731'),(62,62,'齐齐哈尔','Q','2025-05-25 18:21:49.731','2025-05-25 18:21:49.731'),(63,63,'大庆','D','2025-05-25 18:21:49.732','2025-05-25 18:21:49.732'),(64,64,'双鸭山','S','2025-05-25 18:21:49.732','2025-05-25 18:21:49.732'),(65,65,'伊春','Y','2025-05-25 18:21:49.732','2025-05-25 18:21:49.732'),(66,66,'佳木斯','J','2025-05-25 18:21:49.732','2025-05-25 18:21:49.732'),(67,67,'黑河','H','2025-05-25 18:21:49.733','2025-05-25 18:21:49.733'),(68,68,'鹤岗','H','2025-05-25 18:21:49.733','2025-05-25 18:21:49.733'),(69,69,'鸡西','J','2025-05-25 18:21:49.733','2025-05-25 18:21:49.733'),(70,70,'牡丹江','M','2025-05-25 18:21:49.733','2025-05-25 18:21:49.733'),(71,71,'七台河','Q','2025-05-25 18:21:49.733','2025-05-25 18:21:49.733'),(72,72,'绥化','S','2025-05-25 18:21:49.734','2025-05-25 18:21:49.734'),(73,73,'大兴安岭','D','2025-05-25 18:21:49.734','2025-05-25 18:21:49.734'),(74,74,'上海','S','2025-05-25 18:21:49.734','2025-05-25 18:21:49.734'),(75,75,'南京','N','2025-05-25 18:21:49.734','2025-05-25 18:21:49.734'),(76,76,'苏州','S','2025-05-25 18:21:49.735','2025-05-25 18:21:49.735'),(77,77,'扬州','Y','2025-05-25 18:21:49.735','2025-05-25 18:21:49.735'),(78,78,'无锡','W','2025-05-25 18:21:49.735','2025-05-25 18:21:49.735'),(79,79,'徐州','X','2025-05-25 18:21:49.735','2025-05-25 18:21:49.735'),(80,80,'常州','C','2025-05-25 18:21:49.735','2025-05-25 18:21:49.735'),(81,81,'南通','N','2025-05-25 18:21:49.736','2025-05-25 18:21:49.736'),(82,82,'淮安','H','2025-05-25 18:21:49.736','2025-05-25 18:21:49.736'),(83,83,'盐城','Y','2025-05-25 18:21:49.736','2025-05-25 18:21:49.736'),(84,84,'镇江','Z','2025-05-25 18:21:49.736','2025-05-25 18:21:49.736'),(85,85,'泰州','T','2025-05-25 18:21:49.737','2025-05-25 18:21:49.737'),(86,86,'宿迁','S','2025-05-25 18:21:49.737','2025-05-25 18:21:49.737'),(87,87,'连云港','L','2025-05-25 18:21:49.737','2025-05-25 18:21:49.737'),(88,88,'杭州','H','2025-05-25 18:21:49.737','2025-05-25 18:21:49.737'),(89,89,'宁波','N','2025-05-25 18:21:49.738','2025-05-25 18:21:49.738'),(90,90,'温州','W','2025-05-25 18:21:49.738','2025-05-25 18:21:49.738'),(91,91,'嘉兴','J','2025-05-25 18:21:49.738','2025-05-25 18:21:49.738'),(92,92,'湖州','H','2025-05-25 18:21:49.738','2025-05-25 18:21:49.738'),(93,93,'绍兴','S','2025-05-25 18:21:49.739','2025-05-25 18:21:49.739'),(94,94,'金华','J','2025-05-25 18:21:49.739','2025-05-25 18:21:49.739'),(95,95,'衢州','Q','2025-05-25 18:21:49.739','2025-05-25 18:21:49.739'),(96,96,'舟山','Z','2025-05-25 18:21:49.740','2025-05-25 18:21:49.740'),(97,97,'台州','T','2025-05-25 18:21:49.740','2025-05-25 18:21:49.740'),(98,98,'丽水','L','2025-05-25 18:21:49.740','2025-05-25 18:21:49.740'),(99,99,'合肥','H','2025-05-25 18:21:49.740','2025-05-25 18:21:49.740'),(100,100,'芜湖','W','2025-05-25 18:21:49.740','2025-05-25 18:21:49.740'),(101,101,'蚌埠','B','2025-05-25 18:21:49.741','2025-05-25 18:21:49.741'),(102,102,'淮南','H','2025-05-25 18:21:49.741','2025-05-25 18:21:49.741'),(103,103,'马鞍山','M','2025-05-25 18:21:49.741','2025-05-25 18:21:49.741'),(104,104,'淮北','H','2025-05-25 18:21:49.741','2025-05-25 18:21:49.741'),(105,105,'铜陵','T','2025-05-25 18:21:49.741','2025-05-25 18:21:49.741'),(106,106,'滁州','C','2025-05-25 18:21:49.742','2025-05-25 18:21:49.742'),(107,107,'阜阳','F','2025-05-25 18:21:49.742','2025-05-25 18:21:49.742'),(108,108,'亳州','B','2025-05-25 18:21:49.742','2025-05-25 18:21:49.742'),(109,109,'黄山','H','2025-05-25 18:21:49.742','2025-05-25 18:21:49.742'),(110,110,'六安','L','2025-05-25 18:21:49.742','2025-05-25 18:21:49.742'),(111,111,'安庆','A','2025-05-25 18:21:49.743','2025-05-25 18:21:49.743'),(112,112,'宿州','S','2025-05-25 18:21:49.743','2025-05-25 18:21:49.743'),(113,113,'池州','C','2025-05-25 18:21:49.743','2025-05-25 18:21:49.743'),(114,114,'宣城','X','2025-05-25 18:21:49.743','2025-05-25 18:21:49.743'),(115,115,'福州','F','2025-05-25 18:21:49.744','2025-05-25 18:21:49.744'),(116,116,'厦门','X','2025-05-25 18:21:49.744','2025-05-25 18:21:49.744'),(117,117,'泉州','Q','2025-05-25 18:21:49.744','2025-05-25 18:21:49.744'),(118,118,'宁德','N','2025-05-25 18:21:49.744','2025-05-25 18:21:49.744'),(119,119,'龙岩','L','2025-05-25 18:21:49.744','2025-05-25 18:21:49.744'),(120,120,'莆田','P','2025-05-25 18:21:49.745','2025-05-25 18:21:49.745'),(121,121,'三明','S','2025-05-25 18:21:49.745','2025-05-25 18:21:49.745'),(122,122,'漳州','Z','2025-05-25 18:21:49.745','2025-05-25 18:21:49.745'),(123,123,'南平','N','2025-05-25 18:21:49.745','2025-05-25 18:21:49.745'),(124,124,'南昌','N','2025-05-25 18:21:49.746','2025-05-25 18:21:49.746'),(125,125,'九江','J','2025-05-25 18:21:49.746','2025-05-25 18:21:49.746'),(126,126,'抚州','F','2025-05-25 18:21:49.746','2025-05-25 18:21:49.746'),(127,127,'上饶','S','2025-05-25 18:21:49.747','2025-05-25 18:21:49.747'),(128,128,'赣州','G','2025-05-25 18:21:49.747','2025-05-25 18:21:49.747'),(129,129,'景德镇','J','2025-05-25 18:21:49.747','2025-05-25 18:21:49.747'),(130,130,'萍乡','P','2025-05-25 18:21:49.747','2025-05-25 18:21:49.747'),(131,131,'新余','X','2025-05-25 18:21:49.748','2025-05-25 18:21:49.748'),(132,132,'吉安','J','2025-05-25 18:21:49.748','2025-05-25 18:21:49.748'),(133,133,'鹰潭','Y','2025-05-25 18:21:49.748','2025-05-25 18:21:49.748'),(134,134,'宜春','Y','2025-05-25 18:21:49.748','2025-05-25 18:21:49.748'),(135,135,'青岛','Q','2025-05-25 18:21:49.748','2025-05-25 18:21:49.748'),(136,136,'济南','J','2025-05-25 18:21:49.749','2025-05-25 18:21:49.749'),(137,137,'淄博','Z','2025-05-25 18:21:49.749','2025-05-25 18:21:49.749'),(138,138,'枣庄','Z','2025-05-25 18:21:49.749','2025-05-25 18:21:49.749'),(139,139,'东营','D','2025-05-25 18:21:49.749','2025-05-25 18:21:49.749'),(140,140,'烟台','Y','2025-05-25 18:21:49.750','2025-05-25 18:21:49.750'),(141,141,'潍坊','W','2025-05-25 18:21:49.750','2025-05-25 18:21:49.750'),(142,142,'济宁','J','2025-05-25 18:21:49.750','2025-05-25 18:21:49.750'),(143,143,'泰安','T','2025-05-25 18:21:49.750','2025-05-25 18:21:49.750'),(144,144,'威海','W','2025-05-25 18:21:49.751','2025-05-25 18:21:49.751'),(145,145,'日照','R','2025-05-25 18:21:49.751','2025-05-25 18:21:49.751'),(146,146,'临沂','L','2025-05-25 18:21:49.751','2025-05-25 18:21:49.751'),(147,147,'德州','D','2025-05-25 18:21:49.751','2025-05-25 18:21:49.751'),(148,148,'聊城','L','2025-05-25 18:21:49.751','2025-05-25 18:21:49.751'),(149,149,'菏泽','H','2025-05-25 18:21:49.752','2025-05-25 18:21:49.752'),(150,150,'滨州','B','2025-05-25 18:21:49.752','2025-05-25 18:21:49.752'),(151,151,'莱芜','L','2025-05-25 18:21:49.752','2025-05-25 18:21:49.752'),(152,152,'郑州','Z','2025-05-25 18:21:49.752','2025-05-25 18:21:49.752'),(153,153,'开封','K','2025-05-25 18:21:49.753','2025-05-25 18:21:49.753'),(154,154,'洛阳','L','2025-05-25 18:21:49.753','2025-05-25 18:21:49.753'),(155,155,'焦作','J','2025-05-25 18:21:49.753','2025-05-25 18:21:49.753'),(156,156,'信阳','X','2025-05-25 18:21:49.753','2025-05-25 18:21:49.753'),(157,157,'新乡','X','2025-05-25 18:21:49.754','2025-05-25 18:21:49.754'),(158,158,'安阳','A','2025-05-25 18:21:49.754','2025-05-25 18:21:49.754'),(159,159,'漯河','L','2025-05-25 18:21:49.754','2025-05-25 18:21:49.754'),(160,160,'南阳','N','2025-05-25 18:21:49.754','2025-05-25 18:21:49.754'),(161,161,'商丘','S','2025-05-25 18:21:49.755','2025-05-25 18:21:49.755'),(162,162,'驻马店','Z','2025-05-25 18:21:49.755','2025-05-25 18:21:49.755'),(163,163,'济源','J','2025-05-25 18:21:49.755','2025-05-25 18:21:49.755'),(164,164,'鹤壁','H','2025-05-25 18:21:49.755','2025-05-25 18:21:49.755'),(165,165,'濮阳','P','2025-05-25 18:21:49.756','2025-05-25 18:21:49.756'),(166,166,'平顶山','P','2025-05-25 18:21:49.756','2025-05-25 18:21:49.756'),(167,167,'三门峡','S','2025-05-25 18:21:49.756','2025-05-25 18:21:49.756'),(168,168,'许昌','X','2025-05-25 18:21:49.756','2025-05-25 18:21:49.756'),(169,169,'周口','Z','2025-05-25 18:21:49.757','2025-05-25 18:21:49.757'),(170,170,'武汉','W','2025-05-25 18:21:49.757','2025-05-25 18:21:49.757'),(171,171,'宜昌','Y','2025-05-25 18:21:49.757','2025-05-25 18:21:49.757'),(172,172,'襄阳','X','2025-05-25 18:21:49.757','2025-05-25 18:21:49.757'),(173,173,'黄石','H','2025-05-25 18:21:49.758','2025-05-25 18:21:49.758'),(174,174,'黄冈','H','2025-05-25 18:21:49.758','2025-05-25 18:21:49.758'),(175,175,'荆州','J','2025-05-25 18:21:49.758','2025-05-25 18:21:49.758'),(176,176,'咸宁','X','2025-05-25 18:21:49.758','2025-05-25 18:21:49.758'),(177,177,'孝感','X','2025-05-25 18:21:49.758','2025-05-25 18:21:49.758'),(178,178,'随州','S','2025-05-25 18:21:49.759','2025-05-25 18:21:49.759'),(179,179,'仙桃','X','2025-05-25 18:21:49.759','2025-05-25 18:21:49.759'),(180,180,'潜江','Q','2025-05-25 18:21:49.759','2025-05-25 18:21:49.759'),(181,181,'鄂州','E','2025-05-25 18:21:49.759','2025-05-25 18:21:49.759'),(182,182,'荆门','J','2025-05-25 18:21:49.759','2025-05-25 18:21:49.759'),(183,183,'天门','T','2025-05-25 18:21:49.760','2025-05-25 18:21:49.760'),(184,184,'恩施','E','2025-05-25 18:21:49.760','2025-05-25 18:21:49.760'),(185,185,'十堰','S','2025-05-25 18:21:49.760','2025-05-25 18:21:49.760'),(186,186,'神农架','S','2025-05-25 18:21:49.760','2025-05-25 18:21:49.760'),(187,187,'株洲','Z','2025-05-25 18:21:49.760','2025-05-25 18:21:49.760'),(188,188,'衡阳','H','2025-05-25 18:21:49.761','2025-05-25 18:21:49.761'),(189,189,'岳阳','Y','2025-05-25 18:21:49.761','2025-05-25 18:21:49.761'),(190,190,'常德','C','2025-05-25 18:21:49.761','2025-05-25 18:21:49.761'),(191,191,'长沙','C','2025-05-25 18:21:49.762','2025-05-25 18:21:49.762'),(192,192,'娄底','L','2025-05-25 18:21:49.762','2025-05-25 18:21:49.762'),(193,193,'湘潭','X','2025-05-25 18:21:49.762','2025-05-25 18:21:49.762'),(194,194,'郴州','C','2025-05-25 18:21:49.763','2025-05-25 18:21:49.763'),(195,195,'邵阳','S','2025-05-25 18:21:49.763','2025-05-25 18:21:49.763'),(196,196,'益阳','Y','2025-05-25 18:21:49.763','2025-05-25 18:21:49.763'),(197,197,'怀化','H','2025-05-25 18:21:49.763','2025-05-25 18:21:49.763'),(198,198,'湘西','X','2025-05-25 18:21:49.764','2025-05-25 18:21:49.764'),(199,199,'永州','Y','2025-05-25 18:21:49.764','2025-05-25 18:21:49.764'),(200,200,'张家界','Z','2025-05-25 18:21:49.764','2025-05-25 18:21:49.764'),(201,201,'广州','G','2025-05-25 18:21:49.764','2025-05-25 18:21:49.764'),(202,202,'深圳','S','2025-05-25 18:21:49.765','2025-05-25 18:21:49.765'),(203,203,'珠海','Z','2025-05-25 18:21:49.765','2025-05-25 18:21:49.765'),(204,204,'汕头','S','2025-05-25 18:21:49.765','2025-05-25 18:21:49.765'),(205,205,'佛山','F','2025-05-25 18:21:49.765','2025-05-25 18:21:49.765'),(206,206,'江门','J','2025-05-25 18:21:49.765','2025-05-25 18:21:49.765'),(207,207,'湛江','Z','2025-05-25 18:21:49.766','2025-05-25 18:21:49.766'),(208,208,'茂名','M','2025-05-25 18:21:49.766','2025-05-25 18:21:49.766'),(209,209,'惠州','H','2025-05-25 18:21:49.766','2025-05-25 18:21:49.766'),(210,210,'阳江','Y','2025-05-25 18:21:49.766','2025-05-25 18:21:49.766'),(211,211,'清远','Q','2025-05-25 18:21:49.766','2025-05-25 18:21:49.766'),(212,212,'东莞','D','2025-05-25 18:21:49.767','2025-05-25 18:21:49.767'),(213,213,'中山','Z','2025-05-25 18:21:49.767','2025-05-25 18:21:49.767'),(214,214,'肇庆','Z','2025-05-25 18:21:49.767','2025-05-25 18:21:49.767'),(215,215,'河源','H','2025-05-25 18:21:49.767','2025-05-25 18:21:49.767'),(216,216,'潮州','C','2025-05-25 18:21:49.768','2025-05-25 18:21:49.768'),(217,217,'揭阳','J','2025-05-25 18:21:49.768','2025-05-25 18:21:49.768'),(218,218,'韶关','S','2025-05-25 18:21:49.768','2025-05-25 18:21:49.768'),(219,219,'梅州','M','2025-05-25 18:21:49.769','2025-05-25 18:21:49.769'),(220,220,'汕尾','S','2025-05-25 18:21:49.769','2025-05-25 18:21:49.769'),(221,221,'云浮','Y','2025-05-25 18:21:49.769','2025-05-25 18:21:49.769'),(222,222,'南宁','N','2025-05-25 18:21:49.769','2025-05-25 18:21:49.769'),(223,223,'桂林','G','2025-05-25 18:21:49.770','2025-05-25 18:21:49.770'),(224,224,'北海','B','2025-05-25 18:21:49.770','2025-05-25 18:21:49.770'),(225,225,'梧州','W','2025-05-25 18:21:49.770','2025-05-25 18:21:49.770'),(226,226,'钦州','Q','2025-05-25 18:21:49.771','2025-05-25 18:21:49.771'),(227,227,'柳州','L','2025-05-25 18:21:49.771','2025-05-25 18:21:49.771'),(228,228,'贵港','G','2025-05-25 18:21:49.771','2025-05-25 18:21:49.771'),(229,229,'玉林','Y','2025-05-25 18:21:49.771','2025-05-25 18:21:49.771'),(230,230,'百色','B','2025-05-25 18:21:49.772','2025-05-25 18:21:49.772'),(231,231,'贺州','H','2025-05-25 18:21:49.772','2025-05-25 18:21:49.772'),(232,232,'崇左','C','2025-05-25 18:21:49.772','2025-05-25 18:21:49.772'),(233,233,'防城港','F','2025-05-25 18:21:49.772','2025-05-25 18:21:49.772'),(234,234,'河池','H','2025-05-25 18:21:49.773','2025-05-25 18:21:49.773'),(235,235,'来宾','L','2025-05-25 18:21:49.773','2025-05-25 18:21:49.773'),(236,236,'海口','H','2025-05-25 18:21:49.773','2025-05-25 18:21:49.773'),(237,237,'三亚','S','2025-05-25 18:21:49.773','2025-05-25 18:21:49.773'),(238,238,'儋州','D','2025-05-25 18:21:49.773','2025-05-25 18:21:49.773'),(239,239,'琼海','Q','2025-05-25 18:21:49.774','2025-05-25 18:21:49.774'),(240,240,'澄迈','C','2025-05-25 18:21:49.774','2025-05-25 18:21:49.774'),(241,241,'东方','D','2025-05-25 18:21:49.774','2025-05-25 18:21:49.774'),(242,242,'陵水','L','2025-05-25 18:21:49.774','2025-05-25 18:21:49.774'),(243,243,'屯昌','T','2025-05-25 18:21:49.775','2025-05-25 18:21:49.775'),(244,244,'万宁','W','2025-05-25 18:21:49.775','2025-05-25 18:21:49.775'),(245,245,'乐东','L','2025-05-25 18:21:49.775','2025-05-25 18:21:49.775'),(246,246,'文昌','W','2025-05-25 18:21:49.775','2025-05-25 18:21:49.775'),(247,247,'临高县','L','2025-05-25 18:21:49.776','2025-05-25 18:21:49.776'),(248,248,'琼中','Q','2025-05-25 18:21:49.776','2025-05-25 18:21:49.776'),(249,249,'白沙','B','2025-05-25 18:21:49.776','2025-05-25 18:21:49.776'),(250,250,'保亭','B','2025-05-25 18:21:49.776','2025-05-25 18:21:49.776'),(251,251,'昌江','C','2025-05-25 18:21:49.776','2025-05-25 18:21:49.776'),(252,252,'定安','D','2025-05-25 18:21:49.776','2025-05-25 18:21:49.776'),(253,253,'五指山','W','2025-05-25 18:21:49.777','2025-05-25 18:21:49.777'),(254,254,'重庆','C','2025-05-25 18:21:49.777','2025-05-25 18:21:49.777'),(255,255,'成都','C','2025-05-25 18:21:49.777','2025-05-25 18:21:49.777'),(256,256,'绵阳','M','2025-05-25 18:21:49.778','2025-05-25 18:21:49.778'),(257,257,'南充','N','2025-05-25 18:21:49.778','2025-05-25 18:21:49.778'),(258,258,'乐山','L','2025-05-25 18:21:49.778','2025-05-25 18:21:49.778'),(259,259,'宜宾','Y','2025-05-25 18:21:49.778','2025-05-25 18:21:49.778'),(260,260,'达州','D','2025-05-25 18:21:49.785','2025-05-25 18:21:49.785'),(261,261,'资阳','Z','2025-05-25 18:21:49.785','2025-05-25 18:21:49.785'),(262,262,'德阳','D','2025-05-25 18:21:49.786','2025-05-25 18:21:49.786'),(263,263,'遂宁','S','2025-05-25 18:21:49.786','2025-05-25 18:21:49.786'),(264,264,'自贡','Z','2025-05-25 18:21:49.786','2025-05-25 18:21:49.786'),(265,265,'广元','G','2025-05-25 18:21:49.786','2025-05-25 18:21:49.786'),(266,266,'眉山','M','2025-05-25 18:21:49.787','2025-05-25 18:21:49.787'),(267,267,'攀枝花','P','2025-05-25 18:21:49.787','2025-05-25 18:21:49.787'),(268,268,'雅安','Y','2025-05-25 18:21:49.787','2025-05-25 18:21:49.787'),(269,269,'阿坝','A','2025-05-25 18:21:49.787','2025-05-25 18:21:49.787'),(270,270,'巴中','B','2025-05-25 18:21:49.788','2025-05-25 18:21:49.788'),(271,271,'广安','G','2025-05-25 18:21:49.788','2025-05-25 18:21:49.788'),(272,272,'甘孜','G','2025-05-25 18:21:49.788','2025-05-25 18:21:49.788'),(273,273,'泸州','L','2025-05-25 18:21:49.788','2025-05-25 18:21:49.788'),(274,274,'凉山','L','2025-05-25 18:21:49.789','2025-05-25 18:21:49.789'),(275,275,'内江','N','2025-05-25 18:21:49.789','2025-05-25 18:21:49.789'),(276,276,'贵阳','G','2025-05-25 18:21:49.789','2025-05-25 18:21:49.789'),(277,277,'遵义','Z','2025-05-25 18:21:49.789','2025-05-25 18:21:49.789'),(278,278,'铜仁','T','2025-05-25 18:21:49.790','2025-05-25 18:21:49.790'),(279,279,'六盘水','L','2025-05-25 18:21:49.790','2025-05-25 18:21:49.790'),(280,280,'安顺','A','2025-05-25 18:21:49.790','2025-05-25 18:21:49.790'),(281,281,'毕节','B','2025-05-25 18:21:49.790','2025-05-25 18:21:49.790'),(282,282,'黔南','Q','2025-05-25 18:21:49.791','2025-05-25 18:21:49.791'),(283,283,'黔西南','Q','2025-05-25 18:21:49.791','2025-05-25 18:21:49.791'),(284,284,'黔东南','Q','2025-05-25 18:21:49.791','2025-05-25 18:21:49.791'),(285,285,'昆明','K','2025-05-25 18:21:49.791','2025-05-25 18:21:49.791'),(286,286,'大理','D','2025-05-25 18:21:49.792','2025-05-25 18:21:49.792'),(287,287,'红河','H','2025-05-25 18:21:49.792','2025-05-25 18:21:49.792'),(288,288,'玉溪','Y','2025-05-25 18:21:49.792','2025-05-25 18:21:49.792'),(289,289,'普洱','P','2025-05-25 18:21:49.792','2025-05-25 18:21:49.792'),(290,290,'曲靖','Q','2025-05-25 18:21:49.793','2025-05-25 18:21:49.793'),(291,291,'保山','B','2025-05-25 18:21:49.793','2025-05-25 18:21:49.793'),(292,292,'临沧','L','2025-05-25 18:21:49.793','2025-05-25 18:21:49.793'),(293,293,'楚雄','C','2025-05-25 18:21:49.794','2025-05-25 18:21:49.794'),(294,294,'德宏','D','2025-05-25 18:21:49.794','2025-05-25 18:21:49.794'),(295,295,'迪庆','D','2025-05-25 18:21:49.794','2025-05-25 18:21:49.794'),(296,296,'丽江','L','2025-05-25 18:21:49.794','2025-05-25 18:21:49.794'),(297,297,'文山','W','2025-05-25 18:21:49.795','2025-05-25 18:21:49.795'),(298,298,'西双版纳','X','2025-05-25 18:21:49.795','2025-05-25 18:21:49.795'),(299,299,'昭通','Z','2025-05-25 18:21:49.795','2025-05-25 18:21:49.795'),(300,300,'怒江','N','2025-05-25 18:21:49.795','2025-05-25 18:21:49.795'),(301,301,'拉萨','L','2025-05-25 18:21:49.796','2025-05-25 18:21:49.796'),(302,302,'日喀则','R','2025-05-25 18:21:49.796','2025-05-25 18:21:49.796'),(303,303,'那曲','N','2025-05-25 18:21:49.796','2025-05-25 18:21:49.796'),(304,304,'阿里','A','2025-05-25 18:21:49.797','2025-05-25 18:21:49.797'),(305,305,'林芝','L','2025-05-25 18:21:49.797','2025-05-25 18:21:49.797'),(306,306,'昌都','C','2025-05-25 18:21:49.797','2025-05-25 18:21:49.797'),(307,307,'山南','S','2025-05-25 18:21:49.798','2025-05-25 18:21:49.798'),(308,308,'西安','X','2025-05-25 18:21:49.798','2025-05-25 18:21:49.798'),(309,309,'咸阳','X','2025-05-25 18:21:49.798','2025-05-25 18:21:49.798'),(310,310,'宝鸡','B','2025-05-25 18:21:49.798','2025-05-25 18:21:49.798'),(311,311,'汉中','H','2025-05-25 18:21:49.799','2025-05-25 18:21:49.799'),(312,312,'安康','A','2025-05-25 18:21:49.799','2025-05-25 18:21:49.799'),(313,313,'延安','Y','2025-05-25 18:21:49.799','2025-05-25 18:21:49.799'),(314,314,'商洛','S','2025-05-25 18:21:49.800','2025-05-25 18:21:49.800'),(315,315,'铜川','T','2025-05-25 18:21:49.800','2025-05-25 18:21:49.800'),(316,316,'渭南','W','2025-05-25 18:21:49.800','2025-05-25 18:21:49.800'),(317,317,'榆林','Y','2025-05-25 18:21:49.801','2025-05-25 18:21:49.801'),(318,318,'兰州','L','2025-05-25 18:21:49.801','2025-05-25 18:21:49.801'),(319,319,'天水','T','2025-05-25 18:21:49.801','2025-05-25 18:21:49.801'),(320,320,'庆阳','Q','2025-05-25 18:21:49.801','2025-05-25 18:21:49.801'),(321,321,'酒泉','J','2025-05-25 18:21:49.802','2025-05-25 18:21:49.802'),(322,322,'嘉峪关','J','2025-05-25 18:21:49.802','2025-05-25 18:21:49.802'),(323,323,'张掖','Z','2025-05-25 18:21:49.802','2025-05-25 18:21:49.802'),(324,324,'陇南','L','2025-05-25 18:21:49.803','2025-05-25 18:21:49.803'),(325,325,'白银','B','2025-05-25 18:21:49.803','2025-05-25 18:21:49.803'),(326,326,'定西','D','2025-05-25 18:21:49.804','2025-05-25 18:21:49.804'),(327,327,'甘南','G','2025-05-25 18:21:49.804','2025-05-25 18:21:49.804'),(328,328,'金昌','J','2025-05-25 18:21:49.804','2025-05-25 18:21:49.804'),(329,329,'临夏','L','2025-05-25 18:21:49.804','2025-05-25 18:21:49.804'),(330,330,'平凉','P','2025-05-25 18:21:49.805','2025-05-25 18:21:49.805'),(331,331,'武威','W','2025-05-25 18:21:49.805','2025-05-25 18:21:49.805'),(332,332,'西宁','X','2025-05-25 18:21:49.805','2025-05-25 18:21:49.805'),(333,333,'海东','H','2025-05-25 18:21:49.806','2025-05-25 18:21:49.806'),(334,334,'海西','H','2025-05-25 18:21:49.806','2025-05-25 18:21:49.806'),(335,335,'海南州','H','2025-05-25 18:21:49.806','2025-05-25 18:21:49.806'),(336,336,'海北','H','2025-05-25 18:21:49.807','2025-05-25 18:21:49.807'),(337,337,'黄南','H','2025-05-25 18:21:49.807','2025-05-25 18:21:49.807'),(338,338,'玉树','Y','2025-05-25 18:21:49.807','2025-05-25 18:21:49.807'),(339,339,'银川','Y','2025-05-25 18:21:49.808','2025-05-25 18:21:49.808'),(340,340,'石嘴山','S','2025-05-25 18:21:49.808','2025-05-25 18:21:49.808'),(341,341,'固原','G','2025-05-25 18:21:49.808','2025-05-25 18:21:49.808'),(342,342,'中卫','Z','2025-05-25 18:21:49.809','2025-05-25 18:21:49.809'),(343,343,'吴忠','W','2025-05-25 18:21:49.809','2025-05-25 18:21:49.809'),(344,344,'乌鲁木齐','W','2025-05-25 18:21:49.809','2025-05-25 18:21:49.809'),(345,345,'阿克苏','A','2025-05-25 18:21:49.810','2025-05-25 18:21:49.810'),(346,346,'博尔塔拉','B','2025-05-25 18:21:49.810','2025-05-25 18:21:49.810'),(347,347,'巴音郭楞','B','2025-05-25 18:21:49.810','2025-05-25 18:21:49.810'),(348,348,'昌吉','C','2025-05-25 18:21:49.810','2025-05-25 18:21:49.810'),(349,349,'哈密','H','2025-05-25 18:21:49.811','2025-05-25 18:21:49.811'),(350,350,'克拉玛依','K','2025-05-25 18:21:49.811','2025-05-25 18:21:49.811'),(351,351,'喀什','K','2025-05-25 18:21:49.811','2025-05-25 18:21:49.811'),(352,352,'石河子','S','2025-05-25 18:21:49.811','2025-05-25 18:21:49.811'),(353,353,'吐鲁番','T','2025-05-25 18:21:49.811','2025-05-25 18:21:49.811'),(354,354,'塔城','T','2025-05-25 18:21:49.812','2025-05-25 18:21:49.812'),(355,355,'伊犁','Y','2025-05-25 18:21:49.812','2025-05-25 18:21:49.812'),(356,356,'和田','H','2025-05-25 18:21:49.812','2025-05-25 18:21:49.812'),(357,357,'阿勒泰地区','A','2025-05-25 18:21:49.812','2025-05-25 18:21:49.812'),(358,358,'昆玉','K','2025-05-25 18:21:49.813','2025-05-25 18:21:49.813'),(359,359,'克孜勒苏柯尔克孜','K','2025-05-25 18:21:49.813','2025-05-25 18:21:49.813'),(360,360,'五家渠','W','2025-05-25 18:21:49.813','2025-05-25 18:21:49.813'),(361,362,'图木舒克','T','2025-05-25 18:21:49.813','2025-05-25 18:21:49.813'),(362,361,'香港','X','2025-05-25 18:21:49.814','2025-05-25 18:21:49.814');
/*!40000 ALTER TABLE `cities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `distributor_api_logs`
--

DROP TABLE IF EXISTS `distributor_api_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `distributor_api_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `distributor_id` bigint unsigned DEFAULT NULL,
  `api_endpoint` longtext,
  `method` longtext,
  `request_body` longtext,
  `response_code` bigint DEFAULT NULL,
  `response_body` longtext,
  `ip_address` longtext,
  `user_agent` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=54 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `distributor_api_logs`
--

LOCK TABLES `distributor_api_logs` WRITE;
/*!40000 ALTER TABLE `distributor_api_logs` DISABLE KEYS */;
INSERT INTO `distributor_api_logs` VALUES (1,'2025-05-25 12:55:57.589',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(2,'2025-05-25 13:05:47.663',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(3,'2025-05-25 13:06:50.474',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(4,'2025-05-25 13:07:02.318',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(5,'2025-05-25 13:07:19.703',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(6,'2025-05-25 13:09:04.017',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(7,'2025-05-25 13:09:10.550',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(8,'2025-05-25 13:21:42.866',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(9,'2025-05-25 13:26:19.660',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(10,'2025-05-25 13:28:32.276',1,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(11,'2025-05-25 13:28:54.814',1,'/api/v1/distributor/menu','GET','',0,'','::1','PostmanRuntime/7.44.0'),(12,'2025-05-25 13:29:29.776',1,'/api/v1/distributor/goods','GET','',0,'','::1','PostmanRuntime/7.44.0'),(13,'2025-05-25 13:42:40.022',1,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(14,'2025-05-25 13:53:20.060',1,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(15,'2025-05-25 13:54:22.114',1,'/api/v1/distributor/order/DD1748151760023d2d97561','GET','',0,'','::1','PostmanRuntime/7.44.0'),(16,'2025-05-25 18:11:50.988',3,'/api/v1/distributor/cities','GET','',0,'','::1','PostmanRuntime/7.44.0'),(17,'2025-05-25 18:12:28.361',3,'/api/v1/distributor/cities/sync','POST','',0,'','::1','PostmanRuntime/7.44.0'),(18,'2025-05-25 18:12:37.223',3,'/api/v1/distributor/cities/sync','POST','',0,'','::1','PostmanRuntime/7.44.0'),(19,'2025-05-25 18:21:49.406',3,'/api/v1/distributor/cities/sync','POST','',0,'','::1','PostmanRuntime/7.44.0'),(20,'2025-05-25 18:21:55.412',3,'/api/v1/distributor/cities','GET','',0,'','::1','PostmanRuntime/7.44.0'),(21,'2025-05-25 20:25:03.348',3,'/api/v1/distributor/stores','GET','',0,'','::1','PostmanRuntime/7.44.0'),(22,'2025-05-25 20:29:44.973',3,'/api/v1/distributor/cities','GET','',0,'','::1','PostmanRuntime/7.44.0'),(23,'2025-05-28 22:32:26.544',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(24,'2025-05-28 22:37:31.210',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(25,'2025-05-28 22:38:50.015',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(26,'2025-05-28 22:41:24.298',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(27,'2025-05-28 22:50:55.271',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(28,'2025-05-28 22:55:57.205',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(29,'2025-05-28 22:56:03.830',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(30,'2025-05-28 22:59:46.367',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(31,'2025-05-28 22:59:53.970',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(32,'2025-05-28 23:01:34.919',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(33,'2025-05-28 23:15:11.744',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(34,'2025-05-28 23:20:31.055',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(35,'2025-05-28 23:23:43.536',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(36,'2025-05-28 23:23:55.237',3,'/api/v1/distributor/order/DD174844582354738625164','GET','',0,'','::1','PostmanRuntime/7.44.0'),(37,'2025-05-28 23:24:13.126',3,'/api/v1/distributor/order/DD174844582354738625164','GET','',0,'','::1','PostmanRuntime/7.44.0'),(38,'2025-05-28 23:24:13.941',3,'/api/v1/distributor/order/DD174844582354738625164','GET','',0,'','::1','PostmanRuntime/7.44.0'),(39,'2025-05-28 23:25:52.696',3,'/api/v1/distributor/order/DD174844582354738625164','GET','',0,'','::1','PostmanRuntime/7.44.0'),(40,'2025-05-28 23:26:42.415',3,'/api/v1/distributor/order/DD174844582354738625164','GET','',0,'','::1','PostmanRuntime/7.44.0'),(41,'2025-05-28 23:31:01.196',3,'/api/v1/distributor/order/DD174844582354738625164','GET','',0,'','::1','PostmanRuntime/7.44.0'),(42,'2025-05-28 23:38:22.472',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(43,'2025-05-29 21:34:24.027',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(44,'2025-05-29 21:43:46.728',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(45,'2025-05-29 21:44:18.518',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(46,'2025-05-29 21:48:14.084',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(47,'2025-05-29 21:50:40.579',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(48,'2025-05-29 21:50:54.791',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(49,'2025-05-29 21:50:58.418',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(50,'2025-05-29 21:51:33.177',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(51,'2025-05-29 21:56:01.726',3,'/api/v1/distributor/order','POST','',0,'','::1','PostmanRuntime/7.44.0'),(52,'2025-05-29 21:56:21.449',3,'/api/v1/distributor/order/DD1748526961752133ad737','GET','',0,'','::1','PostmanRuntime/7.44.0'),(53,'2025-05-29 22:12:30.834',3,'/api/v1/distributor/order/DD1748526961752133ad737','GET','',0,'','::1','PostmanRuntime/7.44.0');
/*!40000 ALTER TABLE `distributor_api_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `distributors`
--

DROP TABLE IF EXISTS `distributors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `distributors` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `company_name` varchar(200) DEFAULT NULL,
  `contact_name` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `api_key` varchar(64) NOT NULL,
  `api_secret` varchar(64) DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `balance` double DEFAULT NULL,
  `frozen_amount` double DEFAULT NULL,
  `credit_limit` double DEFAULT NULL,
  `callback_url` varchar(500) DEFAULT NULL,
  `warning_balance` double DEFAULT NULL,
  `warning_enabled` tinyint(1) DEFAULT NULL,
  `warning_email` varchar(100) DEFAULT NULL,
  `warning_webhook` varchar(500) DEFAULT NULL,
  `daily_order_limit` bigint DEFAULT NULL,
  `monthly_order_limit` bigint DEFAULT NULL,
  `total_orders` bigint DEFAULT NULL,
  `total_amount` double DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_distributors_api_key` (`api_key`),
  UNIQUE KEY `idx_distributors_email` (`email`),
  KEY `idx_distributors_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `distributors`
--

LOCK TABLES `distributors` WRITE;
/*!40000 ALTER TABLE `distributors` DISABLE KEYS */;
INSERT INTO `distributors` VALUES (1,'2025-05-24 08:57:16.756','2025-05-25 13:42:40.901',NULL,'示例分销商','示例公司','张三','13800138000','demo@distributor.com','a1a71ded-50d9-4edb-8fa4-28995e9e0acc','8d60a52e-52f9-4335-b400-094b5390009d',1,990.9,0,10000,'',0,0,'','',100,3000,0,0,NULL),(2,'2025-05-25 16:17:30.337','2025-05-25 16:17:30.337',NULL,'新分销商','测试公司','联系人','13800138000','test@example.com','e54331ea-aeca-4d13-adb1-b70a39b9a9ac','8bce70c6-b08e-4f65-8d07-774e598357c3',1,0,0,10000,'',0,0,'','',100,3000,0,0,'$2a$10$BmqY4LN/x62JGIiGh/vByufRNwnVPhNw3rVBuBRTjdDxA77SeKKt.'),(3,'2025-05-25 18:11:15.531','2025-05-29 23:36:53.464',NULL,'测试分销商','测试公司','李四','13900139000','test@distributor.com','dfee0899-a324-43b0-9136-1ca96876fb84','0cb6b4d3-0aad-4ca3-8019-f8ad7c82b0d0',1,131.35,0,10000,'https://callback.example.com',0,0,'','',100,3000,0,0,'$2a$10$6yv/muIYuZU93EA0mAC/3.D.RtkYlFPfqxdKSIVdyVYr0B.5ixvxm');
/*!40000 ALTER TABLE `distributors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `luckin_prices`
--

DROP TABLE IF EXISTS `luckin_prices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `luckin_prices` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `price_id` varchar(50) NOT NULL,
  `price_value` double NOT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  `product_codes` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_luckin_prices_price_code` (`price_id`),
  KEY `fk_luckin_prices_creator` (`created_by`),
  CONSTRAINT `fk_luckin_prices_creator` FOREIGN KEY (`created_by`) REFERENCES `admins` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `luckin_prices`
--

LOCK TABLES `luckin_prices` WRITE;
/*!40000 ALTER TABLE `luckin_prices` DISABLE KEYS */;
INSERT INTO `luckin_prices` VALUES (1,'6',9.1,1,'2025-05-24 08:57:16.767','2025-05-24 08:57:16.767',1,NULL),(9,'10',15.5,1,'2025-05-25 15:38:16.197','2025-05-25 15:38:16.197',1,NULL),(10,'15',20,1,'2025-05-25 15:38:16.206','2025-05-25 15:38:16.206',1,NULL),(11,'20',25,1,'2025-05-25 15:38:16.215','2025-05-25 15:38:16.215',1,NULL);
/*!40000 ALTER TABLE `luckin_prices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `luckin_prices_backup_20250529`
--

DROP TABLE IF EXISTS `luckin_prices_backup_20250529`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `luckin_prices_backup_20250529` (
  `id` bigint NOT NULL DEFAULT '0',
  `price_id` varchar(50) NOT NULL,
  `price_value` double NOT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  `product_codes` text
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `luckin_prices_backup_20250529`
--

LOCK TABLES `luckin_prices_backup_20250529` WRITE;
/*!40000 ALTER TABLE `luckin_prices_backup_20250529` DISABLE KEYS */;
INSERT INTO `luckin_prices_backup_20250529` VALUES (1,'6',9.1,1,'2025-05-24 08:57:16.767','2025-05-24 08:57:16.767',1,NULL),(9,'10',15.5,1,'2025-05-25 15:38:16.197','2025-05-25 15:38:16.197',1,NULL),(10,'15',20,1,'2025-05-25 15:38:16.206','2025-05-25 15:38:16.206',1,NULL),(11,'20',25,1,'2025-05-25 15:38:16.215','2025-05-25 15:38:16.215',1,NULL);
/*!40000 ALTER TABLE `luckin_prices_backup_20250529` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `luckin_products`
--

DROP TABLE IF EXISTS `luckin_products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `luckin_products` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` varchar(50) NOT NULL,
  `name` varchar(200) NOT NULL,
  `description` text,
  `category` varchar(100) DEFAULT NULL,
  `image_url` varchar(500) DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_luckin_products_product_id` (`product_id`),
  KEY `fk_luckin_products_creator` (`created_by`),
  CONSTRAINT `fk_luckin_products_creator` FOREIGN KEY (`created_by`) REFERENCES `admins` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `luckin_products`
--

LOCK TABLES `luckin_products` WRITE;
/*!40000 ALTER TABLE `luckin_products` DISABLE KEYS */;
INSERT INTO `luckin_products` VALUES (1,'2500','标准美式','经典美式咖啡','美式家族','',1,'2025-05-24 08:57:16.807','2025-05-24 08:57:16.807',1),(2,'4500','燕麦拿铁','燕麦奶拿铁咖啡','拿铁系列','',1,'2025-05-24 08:57:16.815','2025-05-24 08:57:16.815',1),(3,'4805','拿铁','经典拿铁咖啡','拿铁系列','',1,'2025-05-24 08:57:16.826','2025-05-24 08:57:16.826',1),(4,'4929','橙C冰茶','橙汁茶饮','茶饮系列','',1,'2025-05-24 08:57:16.837','2025-05-24 08:57:16.837',1);
/*!40000 ALTER TABLE `luckin_products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_no` varchar(50) NOT NULL,
  `put_order_id` varchar(50) DEFAULT NULL,
  `distributor_id` bigint unsigned DEFAULT NULL,
  `card_id` bigint unsigned DEFAULT NULL,
  `card_code` varchar(100) DEFAULT NULL,
  `status` bigint DEFAULT NULL,
  `store_code` varchar(50) DEFAULT NULL,
  `store_name` varchar(200) DEFAULT NULL,
  `store_address` varchar(500) DEFAULT NULL,
  `goods` json DEFAULT NULL,
  `total_amount` double DEFAULT NULL,
  `cost_amount` double DEFAULT NULL,
  `profit_amount` double DEFAULT NULL,
  `take_mode` bigint DEFAULT NULL,
  `take_code` varchar(50) DEFAULT NULL,
  `qr_data` text,
  `phone_number` varchar(20) DEFAULT NULL,
  `callback_url` varchar(500) DEFAULT NULL,
  `callback_status` bigint DEFAULT NULL,
  `callback_time` datetime(3) DEFAULT NULL,
  `luckin_response` text,
  `luckin_price` double DEFAULT NULL,
  `luckin_cost_price` double DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_orders_order_no` (`order_no`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `idx_orders_put_order_id` (`put_order_id`),
  KEY `idx_price_diff` (`luckin_price`,`total_amount`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
INSERT INTO `orders` VALUES (1,'2025-05-28 23:23:44.406','2025-05-30 01:35:16.432',NULL,'DD174844582354738625164','1377427157387337728',3,27,'YQC66PDG',2,'362946','福永意库店','宝安区福永街道福海工业区二区一层9号101','[{\"specs\": null, \"goods_id\": \"4500\", \"quantity\": 1, \"sku_code\": \"\", \"sku_name\": \"热/少甜/IIAC金奖豆/大杯 16oz\", \"goods_name\": \"燕麦拿铁\", \"sale_price\": 9.1, \"goods_image\": \"https://img01.luckincoffeecdn.com/group3/M00/5E/48/CtxgFGZ77feAOrjeAAIvfBrDxSA376.png\", \"original_price\": 29}]',9.1,9.1,0,1,'184','aDcqhYb2AAk=','13800138000','https://your-domain.com/callback',1,'2025-05-30 01:35:16.431','{\"code\":200,\"msg\":\"OK\",\"status\":\"success\",\"data\":{\"biz\":\"lk\",\"putOrderId\":\"1377427157387337728\",\"outId\":\"DD174844582354738625164\",\"status\":\"success\",\"takeCode\":\"184\",\"qrdata\":\"aDcqhYb2AAk=\",\"addtime\":\"2025-05-28 23:23:46\",\"storeCode\":\"362946\",\"storeName\":\"福永意库店\",\"salePrice\":\"9.1\",\"costPrice\":\"9.1\",\"refundAmount\":\"0\",\"totalPlatformCostPrice\":\"9.1\",\"refundAll\":null,\"refundAdd\":null,\"goods\":[{\"biz\":\"lk\",\"goodsId\":\"4500\",\"goodsName\":\"燕麦拿铁\",\"num\":1,\"totalOriginalPrice\":\"29\",\"totalSalePrice\":\"9.1\",\"originalPrice\":\"29\",\"salePrice\":\"9.1\",\"refundNum\":0,\"refundAmount\":\"0\",\"memo\":\"\",\"platformCostPrice\":\"9.1\"}],\"takeInfoList\":[{\"takeinfo\":\"184\",\"qrdata\":\"aDcqhYb2AAk=\",\"lockerCode\":\"\",\"orderPhone\":\"8522\",\"goodsId\":\"4500\",\"putOrderDetailId\":15052821,\"detailJson\":\"[{\\\"name\\\":\\\"燕麦拿铁\\\",\\\"image\\\":\\\"https://img01.luckincoffeecdn.com/group3/M00/5E/48/CtxgFGZ77feAOrjeAAIvfBrDxSA376.png\\\",\\\"quantity\\\":1,\\\"skuName\\\":\\\"热/少甜/IIAC金奖豆/大杯 16oz\\\",\\\"comboItems\\\":[]}]\"}]}}',0,0),(2,'2025-05-29 21:56:04.003','2025-05-30 01:35:16.444',NULL,'DD1748526961752133ad737','1377767482781884416',3,29,'KSGHGK4C',2,'387207','益田假日天地1楼店','宝安区怀德社区广深路福永段77号怀施广场一层119室','[{\"specs\": null, \"goods_id\": \"4929\", \"quantity\": 1, \"sku_code\": \"\", \"sku_name\": \"茉莉花香/冰/标准甜/大杯 16oz\", \"goods_name\": \"橙C冰茶\", \"sale_price\": 9.55, \"goods_image\": \"https://img02.luckincoffeecdn.com/group4/M00/07/AF/CtxwDGfgKSuANmzXAAEwOoon23c675.png\", \"original_price\": 26}]',9.55,9.56,-0.009999999999999787,1,'075','aDhneYShAAk=','13800138000','https://your-domain.com/callback',1,'2025-05-30 01:35:16.444','{\"code\":200,\"msg\":\"OK\",\"status\":\"success\",\"data\":{\"biz\":\"lk\",\"putOrderId\":\"1377767482781884416\",\"outId\":\"DD1748526961752133ad737\",\"status\":\"success\",\"takeCode\":\"075\",\"qrdata\":\"aDhneYShAAk=\",\"addtime\":\"2025-05-29 21:56:06\",\"storeCode\":\"387207\",\"storeName\":\"益田假日天地1楼店\",\"salePrice\":\"9.55\",\"costPrice\":\"9.55\",\"refundAmount\":\"0\",\"totalPlatformCostPrice\":\"9.55\",\"refundAll\":null,\"refundAdd\":null,\"goods\":[{\"biz\":\"lk\",\"goodsId\":\"4929\",\"goodsName\":\"橙C冰茶\",\"num\":1,\"totalOriginalPrice\":\"26\",\"totalSalePrice\":\"9.55\",\"originalPrice\":\"26\",\"salePrice\":\"9.55\",\"refundNum\":0,\"refundAmount\":\"0\",\"memo\":\"\",\"platformCostPrice\":\"9.55\"}],\"takeInfoList\":[{\"takeinfo\":\"075\",\"qrdata\":\"aDhneYShAAk=\",\"lockerCode\":\"\",\"orderPhone\":\"7039\",\"goodsId\":\"4929\",\"putOrderDetailId\":15080272,\"detailJson\":\"[{\\\"name\\\":\\\"橙C冰茶\\\",\\\"image\\\":\\\"https://img02.luckincoffeecdn.com/group4/M00/07/AF/CtxwDGfgKSuANmzXAAEwOoon23c675.png\\\",\\\"quantity\\\":1,\\\"skuName\\\":\\\"茉莉花香/冰/标准甜/大杯 16oz\\\",\\\"comboItems\\\":[]}]\"}]}}',0,0);
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_aliases`
--

DROP TABLE IF EXISTS `product_aliases`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_aliases` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_id` bigint unsigned NOT NULL,
  `alias_name` varchar(200) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_aliases_deleted_at` (`deleted_at`),
  KEY `idx_product_aliases_product_id` (`product_id`),
  KEY `idx_product_aliases_alias_name` (`alias_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_aliases`
--

LOCK TABLES `product_aliases` WRITE;
/*!40000 ALTER TABLE `product_aliases` DISABLE KEYS */;
/*!40000 ALTER TABLE `product_aliases` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_match_logs`
--

DROP TABLE IF EXISTS `product_match_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_match_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `distributor_id` bigint unsigned NOT NULL,
  `request_time` datetime(3) NOT NULL,
  `input_product` varchar(200) DEFAULT NULL,
  `input_specs` json DEFAULT NULL,
  `error_reason` varchar(500) DEFAULT NULL,
  `suggestions` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_match_logs_distributor_id` (`distributor_id`),
  KEY `idx_product_match_logs_request_time` (`request_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_match_logs`
--

LOCK TABLES `product_match_logs` WRITE;
/*!40000 ALTER TABLE `product_match_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `product_match_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_price_mappings`
--

DROP TABLE IF EXISTS `product_price_mappings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_price_mappings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `price_id` bigint DEFAULT NULL,
  `product_code` varchar(50) DEFAULT NULL,
  `sku_code` varchar(50) DEFAULT NULL,
  `priority` bigint DEFAULT '0',
  `status` bigint DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `idx_product_price_mappings_deleted_at` (`deleted_at`),
  KEY `idx_product_price_mappings_price_id` (`price_id`),
  KEY `idx_product_price_mappings_product_code` (`product_code`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_price_mappings`
--

LOCK TABLES `product_price_mappings` WRITE;
/*!40000 ALTER TABLE `product_price_mappings` DISABLE KEYS */;
INSERT INTO `product_price_mappings` VALUES (1,'2025-05-28 22:31:41.499','2025-05-28 22:31:41.499',NULL,0,'2500','',0,1),(2,'2025-05-28 22:31:41.679','2025-05-28 22:31:41.679',NULL,0,'4500','',0,1),(3,'2025-05-28 22:31:41.938','2025-05-28 22:31:41.938',NULL,0,'4805','',0,1),(4,'2025-05-28 22:31:42.085','2025-05-28 22:31:42.085',NULL,0,'4929','',0,1),(5,'2025-05-28 22:56:02.709','2025-05-28 22:56:02.709',NULL,6,'2500','',0,1),(6,'2025-05-28 22:56:02.909','2025-05-28 22:56:02.909',NULL,6,'4500','',0,1),(7,'2025-05-28 22:56:03.177','2025-05-28 22:56:03.177',NULL,6,'4805','',0,1),(8,'2025-05-28 22:56:03.346','2025-05-28 22:56:03.346',NULL,6,'4929','',0,1),(9,'2025-05-29 19:46:24.150','2025-05-29 19:46:24.150',NULL,1,'2500','',0,1),(10,'2025-05-29 19:46:24.370','2025-05-29 19:46:24.370',NULL,1,'4500','',0,1),(11,'2025-05-29 19:46:24.671','2025-05-29 19:46:24.671',NULL,1,'4805','',0,1),(12,'2025-05-29 19:46:24.836','2025-05-29 19:46:24.836',NULL,1,'4929','',0,1);
/*!40000 ALTER TABLE `product_price_mappings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_sku_mappings`
--

DROP TABLE IF EXISTS `product_sku_mappings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_sku_mappings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_id` bigint unsigned NOT NULL,
  `sku_code` varchar(50) NOT NULL,
  `chinese_desc` varchar(200) NOT NULL,
  `specs_code` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_sku_mappings_deleted_at` (`deleted_at`),
  KEY `idx_product_sku_mappings_product_id` (`product_id`),
  KEY `idx_product_sku_mappings_chinese_desc` (`chinese_desc`)
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_sku_mappings`
--

LOCK TABLES `product_sku_mappings` WRITE;
/*!40000 ALTER TABLE `product_sku_mappings` DISABLE KEYS */;
INSERT INTO `product_sku_mappings` VALUES (1,'2025-05-28 21:12:22.759','2025-05-28 22:31:41.244',NULL,1,'SP2277-00223','大杯 16ozIIAC金奖豆冰不另外加糖双份奶','0_0_0_4_0'),(2,'2025-05-28 21:12:22.768','2025-05-28 22:31:41.253',NULL,1,'SP2277-00227','大杯 16ozIIAC金奖豆冰少甜单份奶','0_0_0_1_1'),(3,'2025-05-28 21:12:22.776','2025-05-28 22:31:41.261',NULL,1,'SP2277-00228','大杯 16ozIIAC金奖豆冰微甜单份奶','0_0_0_3_1'),(4,'2025-05-28 21:12:22.785','2025-05-28 22:31:41.268',NULL,1,'SP2277-00192','大杯 16ozIIAC金奖豆热标准甜单份奶','0_0_1_0_1'),(5,'2025-05-28 21:12:22.826','2025-05-28 22:31:41.276',NULL,1,'SP2277-00197','大杯 16ozIIAC金奖豆冰不另外加糖单份奶','0_0_0_4_1'),(6,'2025-05-28 21:12:22.856','2025-05-28 22:31:41.284',NULL,1,'SP2277-00211','大杯 16ozIIAC金奖豆热微甜单份奶','0_0_1_3_1'),(7,'2025-05-28 21:12:22.863','2025-05-28 22:31:41.292',NULL,1,'SP2277-00214','大杯 16ozIIAC金奖豆热少甜单份奶','0_0_1_1_1'),(8,'2025-05-28 21:12:22.872','2025-05-28 22:31:41.299',NULL,1,'SP2277-00185','大杯 16ozIIAC金奖豆热不另外加糖双份奶','0_0_1_4_0'),(9,'2025-05-28 21:12:22.882','2025-05-28 22:31:41.306',NULL,1,'SP2277-00183','大杯 16ozIIAC金奖豆热标准甜双份奶','0_0_1_0_0'),(10,'2025-05-28 21:12:22.889','2025-05-28 22:31:41.313',NULL,1,'SP2277-00248','大杯 16ozIIAC金奖豆热少甜双份奶','0_0_1_1_0'),(11,'2025-05-28 21:12:22.899','2025-05-28 22:31:41.320',NULL,1,'SP2277-00247','大杯 16ozIIAC金奖豆热微甜双份奶','0_0_1_3_0'),(12,'2025-05-28 21:12:22.909','2025-05-28 22:31:41.328',NULL,1,'SP2277-00252','大杯 16ozIIAC金奖豆热少少甜单份奶','0_0_1_2_1'),(13,'2025-05-28 21:12:22.919','2025-05-28 22:31:41.336',NULL,1,'SP2277-00233','大杯 16ozIIAC金奖豆冰少甜无奶','0_0_0_1_2'),(14,'2025-05-28 21:12:22.929','2025-05-28 22:31:41.348',NULL,1,'SP2277-00234','大杯 16ozIIAC金奖豆热少甜无奶','0_0_1_1_2'),(15,'2025-05-28 21:12:22.938','2025-05-28 22:31:41.356',NULL,1,'SP2277-00232','大杯 16ozIIAC金奖豆冰微甜无奶','0_0_0_3_2'),(16,'2025-05-28 21:12:22.952','2025-05-28 22:31:41.365',NULL,1,'SP2277-00235','大杯 16ozIIAC金奖豆热微甜无奶','0_0_1_3_2'),(17,'2025-05-28 21:12:22.963','2025-05-28 22:31:41.373',NULL,1,'SP2277-00241','大杯 16ozIIAC金奖豆热不另外加糖单份奶','0_0_1_4_1'),(18,'2025-05-28 21:12:22.973','2025-05-28 22:31:41.381',NULL,1,'SP2277-00266','大杯 16ozIIAC金奖豆热不另外加糖无奶','0_0_1_4_2'),(19,'2025-05-28 21:12:22.980','2025-05-28 22:31:41.387',NULL,1,'SP2277-00264','大杯 16ozIIAC金奖豆冰少少甜无奶','0_0_0_2_2'),(20,'2025-05-28 21:12:22.989','2025-05-28 22:31:41.395',NULL,1,'SP2277-00271','大杯 16ozIIAC金奖豆热少少甜双份奶','0_0_1_2_0'),(21,'2025-05-28 21:12:22.996','2025-05-28 22:31:41.403',NULL,1,'SP2277-00200','大杯 16ozIIAC金奖豆热少少甜无奶','0_0_1_2_2'),(22,'2025-05-28 21:12:23.003','2025-05-28 22:31:41.411',NULL,1,'SP2277-00168','大杯 16ozIIAC金奖豆冰少少甜双份奶','0_0_0_2_0'),(23,'2025-05-28 21:12:23.012','2025-05-28 22:31:41.418',NULL,1,'SP2277-00286','大杯 16ozIIAC金奖豆冰不另外加糖无奶','0_0_0_4_2'),(24,'2025-05-28 21:12:23.021','2025-05-28 22:31:41.427',NULL,1,'SP2277-00287','大杯 16ozIIAC金奖豆冰少甜双份奶','0_0_0_1_0'),(25,'2025-05-28 21:12:23.030','2025-05-28 22:31:41.435',NULL,1,'SP2277-00203','大杯 16ozIIAC金奖豆热标准甜无奶','0_0_1_0_2'),(26,'2025-05-28 21:12:23.041','2025-05-28 22:31:41.445',NULL,1,'SP2277-00174','大杯 16ozIIAC金奖豆冰标准甜无奶','0_0_0_0_2'),(27,'2025-05-28 21:12:23.050','2025-05-28 22:31:41.458',NULL,1,'SP2277-00296','大杯 16ozIIAC金奖豆冰标准甜单份奶','0_0_0_0_1'),(28,'2025-05-28 21:12:23.058','2025-05-28 22:31:41.467',NULL,1,'SP2277-00293','大杯 16ozIIAC金奖豆冰微甜双份奶','0_0_0_3_0'),(29,'2025-05-28 21:12:23.065','2025-05-28 22:31:41.481',NULL,1,'SP2277-00279','大杯 16ozIIAC金奖豆冰少少甜单份奶','0_0_0_2_1'),(30,'2025-05-28 21:12:23.073','2025-05-28 22:31:41.496',NULL,1,'SP2277-00283','大杯 16ozIIAC金奖豆冰标准甜双份奶','0_0_0_0_0'),(31,'2025-05-28 21:12:23.202','2025-05-28 22:31:41.625',NULL,2,'SP2983-00042','大杯 16ozIIAC金奖豆冰标准甜','0_0_0_0'),(32,'2025-05-28 21:12:23.209','2025-05-28 22:31:41.632',NULL,2,'SP2983-00064','大杯 16ozIIAC金奖豆冰少甜','0_0_0_1'),(33,'2025-05-28 21:12:23.217','2025-05-28 22:31:41.639',NULL,2,'SP2983-00066','大杯 16ozIIAC金奖豆冰不另外加糖','0_0_0_4'),(34,'2025-05-28 21:12:23.223','2025-05-28 22:31:41.644',NULL,2,'SP2983-00048','大杯 16ozIIAC金奖豆冰少少甜','0_0_0_2'),(35,'2025-05-28 21:12:23.230','2025-05-28 22:31:41.651',NULL,2,'SP2983-00053','大杯 16ozIIAC金奖豆冰微甜','0_0_0_3'),(36,'2025-05-28 21:12:23.236','2025-05-28 22:31:41.656',NULL,2,'SP2983-00031','大杯 16ozIIAC金奖豆热不另外加糖','0_0_1_4'),(37,'2025-05-28 21:12:23.242','2025-05-28 22:31:41.661',NULL,2,'SP2983-00056','大杯 16ozIIAC金奖豆热微甜','0_0_1_3'),(38,'2025-05-28 21:12:23.250','2025-05-28 22:31:41.668',NULL,2,'SP2983-00055','大杯 16ozIIAC金奖豆热标准甜','0_0_1_0'),(39,'2025-05-28 21:12:23.259','2025-05-28 22:31:41.673',NULL,2,'SP2983-00038','大杯 16ozIIAC金奖豆热少甜','0_0_1_1'),(40,'2025-05-28 21:12:23.266','2025-05-28 22:31:41.678',NULL,2,'SP2983-00063','大杯 16ozIIAC金奖豆热少少甜','0_0_1_2'),(41,'2025-05-28 21:12:23.397','2025-05-28 22:31:41.806',NULL,3,'SP3225-00050','大杯 16ozIIAC金奖豆热少甜纯牛奶','0_0_1_1_0'),(42,'2025-05-28 21:12:23.405','2025-05-28 22:31:41.814',NULL,3,'SP3225-00052','大杯 16ozIIAC金奖豆热标准甜燕麦奶','0_0_1_0_1'),(43,'2025-05-28 21:12:23.413','2025-05-28 22:31:41.820',NULL,3,'SP3225-00051','大杯 16ozIIAC金奖豆热标准甜纯牛奶','0_0_1_0_0'),(44,'2025-05-28 21:12:23.421','2025-05-28 22:31:41.828',NULL,3,'SP3225-00042','大杯 16ozIIAC金奖豆热少少甜纯牛奶','0_0_1_2_0'),(45,'2025-05-28 21:12:23.430','2025-05-28 22:31:41.835',NULL,3,'SP3225-00044','大杯 16ozIIAC金奖豆热少少甜燕麦奶','0_0_1_2_1'),(46,'2025-05-28 21:12:23.438','2025-05-28 22:31:41.843',NULL,3,'SP3225-00046','大杯 16ozIIAC金奖豆冰少甜燕麦奶','0_0_0_1_1'),(47,'2025-05-28 21:12:23.446','2025-05-28 22:31:41.849',NULL,3,'SP3225-00049','大杯 16ozIIAC金奖豆冰微甜纯牛奶','0_0_0_3_0'),(48,'2025-05-28 21:12:23.453','2025-05-28 22:31:41.856',NULL,3,'SP3225-00048','大杯 16ozIIAC金奖豆冰少甜纯牛奶','0_0_0_1_0'),(49,'2025-05-28 21:12:23.461','2025-05-28 22:31:41.863',NULL,3,'SP3225-00075','大杯 16ozIIAC金奖豆热不另外加糖燕麦奶','0_0_1_4_1'),(50,'2025-05-28 21:12:23.471','2025-05-28 22:31:41.870',NULL,3,'SP3225-00078','大杯 16ozIIAC金奖豆热微甜纯牛奶','0_0_1_3_0'),(51,'2025-05-28 21:12:23.482','2025-05-28 22:31:41.877',NULL,3,'SP3225-00077','大杯 16ozIIAC金奖豆热不另外加糖纯牛奶','0_0_1_4_0'),(52,'2025-05-28 21:12:23.490','2025-05-28 22:31:41.883',NULL,3,'SP3225-00079','大杯 16ozIIAC金奖豆热微甜燕麦奶','0_0_1_3_1'),(53,'2025-05-28 21:12:23.499','2025-05-28 22:31:41.889',NULL,3,'SP3225-00072','大杯 16ozIIAC金奖豆冰不另外加糖燕麦奶','0_0_0_4_1'),(54,'2025-05-28 21:12:23.505','2025-05-28 22:31:41.896',NULL,3,'SP3225-00065','大杯 16ozIIAC金奖豆冰少少甜纯牛奶','0_0_0_2_0'),(55,'2025-05-28 21:12:23.513','2025-05-28 22:31:41.902',NULL,3,'SP3225-00067','大杯 16ozIIAC金奖豆冰不另外加糖纯牛奶','0_0_0_4_0'),(56,'2025-05-28 21:12:23.519','2025-05-28 22:31:41.909',NULL,3,'SP3225-00066','大杯 16ozIIAC金奖豆冰少少甜燕麦奶','0_0_0_2_1'),(57,'2025-05-28 21:12:23.527','2025-05-28 22:31:41.916',NULL,3,'SP3225-00061','大杯 16ozIIAC金奖豆冰标准甜燕麦奶','0_0_0_0_1'),(58,'2025-05-28 21:12:23.533','2025-05-28 22:31:41.922',NULL,3,'SP3225-00060','大杯 16ozIIAC金奖豆冰标准甜纯牛奶','0_0_0_0_0'),(59,'2025-05-28 21:12:23.540','2025-05-28 22:31:41.930',NULL,3,'SP3225-00062','大杯 16ozIIAC金奖豆冰微甜燕麦奶','0_0_0_3_1'),(60,'2025-05-28 21:12:23.546','2025-05-28 22:31:41.937',NULL,3,'SP3225-00058','大杯 16ozIIAC金奖豆热少甜燕麦奶','0_0_1_1_1'),(61,'2025-05-28 21:12:23.681','2025-05-28 22:31:42.061',NULL,4,'SP3349-00006','大杯 16oz冰少少甜茉莉花香','0_0_2_0'),(62,'2025-05-28 21:12:23.688','2025-05-28 22:31:42.069',NULL,4,'SP3349-00008','大杯 16oz冰标准甜茉莉花香','0_0_0_0'),(63,'2025-05-28 21:12:23.695','2025-05-28 22:31:42.074',NULL,4,'SP3349-00007','大杯 16oz冰不另外加糖茉莉花香','0_0_4_0'),(64,'2025-05-28 21:12:23.701','2025-05-28 22:31:42.079',NULL,4,'SP3349-00009','大杯 16oz冰微甜茉莉花香','0_0_3_0'),(65,'2025-05-28 21:12:23.708','2025-05-28 22:31:42.084',NULL,4,'SP3349-00010','大杯 16oz冰少甜茉莉花香','0_0_1_0');
/*!40000 ALTER TABLE `product_sku_mappings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_skus`
--

DROP TABLE IF EXISTS `product_skus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_skus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_id` bigint unsigned DEFAULT NULL,
  `sku_code` varchar(50) NOT NULL,
  `sku_name` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_skus_deleted_at` (`deleted_at`),
  KEY `idx_product_skus_product_id` (`product_id`),
  KEY `idx_product_skus_sku_code` (`sku_code`)
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_skus`
--

LOCK TABLES `product_skus` WRITE;
/*!40000 ALTER TABLE `product_skus` DISABLE KEYS */;
INSERT INTO `product_skus` VALUES (1,'2025-05-28 21:12:22.743','2025-05-28 21:12:22.743',NULL,1,'SP2277-00223','大杯 16oz/IIAC金奖豆/冰/不另外加糖/双份奶'),(2,'2025-05-28 21:12:22.760','2025-05-28 21:12:22.760',NULL,1,'SP2277-00227','大杯 16oz/IIAC金奖豆/冰/少甜/单份奶'),(3,'2025-05-28 21:12:22.768','2025-05-28 21:12:22.768',NULL,1,'SP2277-00228','大杯 16oz/IIAC金奖豆/冰/微甜/单份奶'),(4,'2025-05-28 21:12:22.777','2025-05-28 21:12:22.777',NULL,1,'SP2277-00192','大杯 16oz/IIAC金奖豆/热/标准甜/单份奶'),(5,'2025-05-28 21:12:22.786','2025-05-28 21:12:22.786',NULL,1,'SP2277-00197','大杯 16oz/IIAC金奖豆/冰/不另外加糖/单份奶'),(6,'2025-05-28 21:12:22.849','2025-05-28 21:12:22.849',NULL,1,'SP2277-00211','大杯 16oz/IIAC金奖豆/热/微甜/单份奶'),(7,'2025-05-28 21:12:22.857','2025-05-28 21:12:22.857',NULL,1,'SP2277-00214','大杯 16oz/IIAC金奖豆/热/少甜/单份奶'),(8,'2025-05-28 21:12:22.864','2025-05-28 21:12:22.864',NULL,1,'SP2277-00185','大杯 16oz/IIAC金奖豆/热/不另外加糖/双份奶'),(9,'2025-05-28 21:12:22.872','2025-05-28 21:12:22.872',NULL,1,'SP2277-00183','大杯 16oz/IIAC金奖豆/热/标准甜/双份奶'),(10,'2025-05-28 21:12:22.883','2025-05-28 21:12:22.883',NULL,1,'SP2277-00248','大杯 16oz/IIAC金奖豆/热/少甜/双份奶'),(11,'2025-05-28 21:12:22.890','2025-05-28 21:12:22.890',NULL,1,'SP2277-00247','大杯 16oz/IIAC金奖豆/热/微甜/双份奶'),(12,'2025-05-28 21:12:22.900','2025-05-28 21:12:22.900',NULL,1,'SP2277-00252','大杯 16oz/IIAC金奖豆/热/少少甜/单份奶'),(13,'2025-05-28 21:12:22.910','2025-05-28 21:12:22.910',NULL,1,'SP2277-00233','大杯 16oz/IIAC金奖豆/冰/少甜/无奶'),(14,'2025-05-28 21:12:22.920','2025-05-28 21:12:22.920',NULL,1,'SP2277-00234','大杯 16oz/IIAC金奖豆/热/少甜/无奶'),(15,'2025-05-28 21:12:22.930','2025-05-28 21:12:22.930',NULL,1,'SP2277-00232','大杯 16oz/IIAC金奖豆/冰/微甜/无奶'),(16,'2025-05-28 21:12:22.939','2025-05-28 21:12:22.939',NULL,1,'SP2277-00235','大杯 16oz/IIAC金奖豆/热/微甜/无奶'),(17,'2025-05-28 21:12:22.953','2025-05-28 21:12:22.953',NULL,1,'SP2277-00241','大杯 16oz/IIAC金奖豆/热/不另外加糖/单份奶'),(18,'2025-05-28 21:12:22.964','2025-05-28 21:12:22.964',NULL,1,'SP2277-00266','大杯 16oz/IIAC金奖豆/热/不另外加糖/无奶'),(19,'2025-05-28 21:12:22.974','2025-05-28 21:12:22.974',NULL,1,'SP2277-00264','大杯 16oz/IIAC金奖豆/冰/少少甜/无奶'),(20,'2025-05-28 21:12:22.981','2025-05-28 21:12:22.981',NULL,1,'SP2277-00271','大杯 16oz/IIAC金奖豆/热/少少甜/双份奶'),(21,'2025-05-28 21:12:22.989','2025-05-28 21:12:22.989',NULL,1,'SP2277-00200','大杯 16oz/IIAC金奖豆/热/少少甜/无奶'),(22,'2025-05-28 21:12:22.996','2025-05-28 21:12:22.996',NULL,1,'SP2277-00168','大杯 16oz/IIAC金奖豆/冰/少少甜/双份奶'),(23,'2025-05-28 21:12:23.004','2025-05-28 21:12:23.004',NULL,1,'SP2277-00286','大杯 16oz/IIAC金奖豆/冰/不另外加糖/无奶'),(24,'2025-05-28 21:12:23.013','2025-05-28 21:12:23.013',NULL,1,'SP2277-00287','大杯 16oz/IIAC金奖豆/冰/少甜/双份奶'),(25,'2025-05-28 21:12:23.022','2025-05-28 21:12:23.022',NULL,1,'SP2277-00203','大杯 16oz/IIAC金奖豆/热/标准甜/无奶'),(26,'2025-05-28 21:12:23.031','2025-05-28 21:12:23.031',NULL,1,'SP2277-00174','大杯 16oz/IIAC金奖豆/冰/标准甜/无奶'),(27,'2025-05-28 21:12:23.042','2025-05-28 21:12:23.042',NULL,1,'SP2277-00296','大杯 16oz/IIAC金奖豆/冰/标准甜/单份奶'),(28,'2025-05-28 21:12:23.051','2025-05-28 21:12:23.051',NULL,1,'SP2277-00293','大杯 16oz/IIAC金奖豆/冰/微甜/双份奶'),(29,'2025-05-28 21:12:23.059','2025-05-28 21:12:23.059',NULL,1,'SP2277-00279','大杯 16oz/IIAC金奖豆/冰/少少甜/单份奶'),(30,'2025-05-28 21:12:23.066','2025-05-28 21:12:23.066',NULL,1,'SP2277-00283','大杯 16oz/IIAC金奖豆/冰/标准甜/双份奶'),(31,'2025-05-28 21:12:23.190','2025-05-28 21:12:23.190',NULL,2,'SP2983-00042','大杯 16oz/IIAC金奖豆/冰/标准甜'),(32,'2025-05-28 21:12:23.203','2025-05-28 21:12:23.203',NULL,2,'SP2983-00064','大杯 16oz/IIAC金奖豆/冰/少甜'),(33,'2025-05-28 21:12:23.210','2025-05-28 21:12:23.210',NULL,2,'SP2983-00066','大杯 16oz/IIAC金奖豆/冰/不另外加糖'),(34,'2025-05-28 21:12:23.218','2025-05-28 21:12:23.218',NULL,2,'SP2983-00048','大杯 16oz/IIAC金奖豆/冰/少少甜'),(35,'2025-05-28 21:12:23.224','2025-05-28 21:12:23.224',NULL,2,'SP2983-00053','大杯 16oz/IIAC金奖豆/冰/微甜'),(36,'2025-05-28 21:12:23.230','2025-05-28 21:12:23.230',NULL,2,'SP2983-00031','大杯 16oz/IIAC金奖豆/热/不另外加糖'),(37,'2025-05-28 21:12:23.237','2025-05-28 21:12:23.237',NULL,2,'SP2983-00056','大杯 16oz/IIAC金奖豆/热/微甜'),(38,'2025-05-28 21:12:23.244','2025-05-28 21:12:23.244',NULL,2,'SP2983-00055','大杯 16oz/IIAC金奖豆/热/标准甜'),(39,'2025-05-28 21:12:23.251','2025-05-28 21:12:23.251',NULL,2,'SP2983-00038','大杯 16oz/IIAC金奖豆/热/少甜'),(40,'2025-05-28 21:12:23.259','2025-05-28 21:12:23.259',NULL,2,'SP2983-00063','大杯 16oz/IIAC金奖豆/热/少少甜'),(41,'2025-05-28 21:12:23.386','2025-05-28 21:12:23.386',NULL,3,'SP3225-00050','大杯 16oz/IIAC金奖豆/热/少甜/纯牛奶'),(42,'2025-05-28 21:12:23.398','2025-05-28 21:12:23.398',NULL,3,'SP3225-00052','大杯 16oz/IIAC金奖豆/热/标准甜/燕麦奶'),(43,'2025-05-28 21:12:23.406','2025-05-28 21:12:23.406',NULL,3,'SP3225-00051','大杯 16oz/IIAC金奖豆/热/标准甜/纯牛奶'),(44,'2025-05-28 21:12:23.414','2025-05-28 21:12:23.414',NULL,3,'SP3225-00042','大杯 16oz/IIAC金奖豆/热/少少甜/纯牛奶'),(45,'2025-05-28 21:12:23.422','2025-05-28 21:12:23.422',NULL,3,'SP3225-00044','大杯 16oz/IIAC金奖豆/热/少少甜/燕麦奶'),(46,'2025-05-28 21:12:23.431','2025-05-28 21:12:23.431',NULL,3,'SP3225-00046','大杯 16oz/IIAC金奖豆/冰/少甜/燕麦奶'),(47,'2025-05-28 21:12:23.439','2025-05-28 21:12:23.439',NULL,3,'SP3225-00049','大杯 16oz/IIAC金奖豆/冰/微甜/纯牛奶'),(48,'2025-05-28 21:12:23.447','2025-05-28 21:12:23.447',NULL,3,'SP3225-00048','大杯 16oz/IIAC金奖豆/冰/少甜/纯牛奶'),(49,'2025-05-28 21:12:23.454','2025-05-28 21:12:23.454',NULL,3,'SP3225-00075','大杯 16oz/IIAC金奖豆/热/不另外加糖/燕麦奶'),(50,'2025-05-28 21:12:23.462','2025-05-28 21:12:23.462',NULL,3,'SP3225-00078','大杯 16oz/IIAC金奖豆/热/微甜/纯牛奶'),(51,'2025-05-28 21:12:23.472','2025-05-28 21:12:23.472',NULL,3,'SP3225-00077','大杯 16oz/IIAC金奖豆/热/不另外加糖/纯牛奶'),(52,'2025-05-28 21:12:23.483','2025-05-28 21:12:23.483',NULL,3,'SP3225-00079','大杯 16oz/IIAC金奖豆/热/微甜/燕麦奶'),(53,'2025-05-28 21:12:23.491','2025-05-28 21:12:23.491',NULL,3,'SP3225-00072','大杯 16oz/IIAC金奖豆/冰/不另外加糖/燕麦奶'),(54,'2025-05-28 21:12:23.499','2025-05-28 21:12:23.499',NULL,3,'SP3225-00065','大杯 16oz/IIAC金奖豆/冰/少少甜/纯牛奶'),(55,'2025-05-28 21:12:23.506','2025-05-28 21:12:23.506',NULL,3,'SP3225-00067','大杯 16oz/IIAC金奖豆/冰/不另外加糖/纯牛奶'),(56,'2025-05-28 21:12:23.514','2025-05-28 21:12:23.514',NULL,3,'SP3225-00066','大杯 16oz/IIAC金奖豆/冰/少少甜/燕麦奶'),(57,'2025-05-28 21:12:23.520','2025-05-28 21:12:23.520',NULL,3,'SP3225-00061','大杯 16oz/IIAC金奖豆/冰/标准甜/燕麦奶'),(58,'2025-05-28 21:12:23.528','2025-05-28 21:12:23.528',NULL,3,'SP3225-00060','大杯 16oz/IIAC金奖豆/冰/标准甜/纯牛奶'),(59,'2025-05-28 21:12:23.534','2025-05-28 21:12:23.534',NULL,3,'SP3225-00062','大杯 16oz/IIAC金奖豆/冰/微甜/燕麦奶'),(60,'2025-05-28 21:12:23.541','2025-05-28 21:12:23.541',NULL,3,'SP3225-00058','大杯 16oz/IIAC金奖豆/热/少甜/燕麦奶'),(61,'2025-05-28 21:12:23.673','2025-05-28 21:12:23.673',NULL,4,'SP3349-00006','大杯 16oz/冰/少少甜/茉莉花香'),(62,'2025-05-28 21:12:23.682','2025-05-28 21:12:23.682',NULL,4,'SP3349-00008','大杯 16oz/冰/标准甜/茉莉花香'),(63,'2025-05-28 21:12:23.689','2025-05-28 21:12:23.689',NULL,4,'SP3349-00007','大杯 16oz/冰/不另外加糖/茉莉花香'),(64,'2025-05-28 21:12:23.695','2025-05-28 21:12:23.695',NULL,4,'SP3349-00009','大杯 16oz/冰/微甜/茉莉花香'),(65,'2025-05-28 21:12:23.702','2025-05-28 21:12:23.702',NULL,4,'SP3349-00010','大杯 16oz/冰/少甜/茉莉花香');
/*!40000 ALTER TABLE `product_skus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_spec_configs`
--

DROP TABLE IF EXISTS `product_spec_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_spec_configs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_id` bigint unsigned NOT NULL,
  `spec_type` varchar(50) NOT NULL,
  `chinese_name` varchar(50) NOT NULL,
  `is_required` tinyint(1) DEFAULT '1',
  `default_value` varchar(50) DEFAULT NULL,
  `display_order` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_product_spec_configs_deleted_at` (`deleted_at`),
  KEY `idx_product_spec_configs_product_id` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_spec_configs`
--

LOCK TABLES `product_spec_configs` WRITE;
/*!40000 ALTER TABLE `product_spec_configs` DISABLE KEYS */;
INSERT INTO `product_spec_configs` VALUES (1,'2025-05-28 21:12:22.745','2025-05-28 21:12:22.747',NULL,1,'size','杯型',1,'大杯 16oz',0),(2,'2025-05-28 21:12:22.748','2025-05-28 21:12:22.749',NULL,1,'bean','咖啡豆',1,'IIAC金奖豆',1),(3,'2025-05-28 21:12:22.750','2025-05-28 21:12:22.781',NULL,1,'temperature','温度',1,'热',2),(4,'2025-05-28 21:12:22.753','2025-05-28 21:12:22.756',NULL,1,'sweetness','甜度',1,'不另外加糖',3),(5,'2025-05-28 21:12:22.757','2025-05-28 21:12:22.918',NULL,1,'milk','奶',1,'无奶',4),(6,'2025-05-28 21:12:23.191','2025-05-28 21:12:23.192',NULL,2,'size','杯型',1,'大杯 16oz',0),(7,'2025-05-28 21:12:23.193','2025-05-28 21:12:23.194',NULL,2,'bean','咖啡豆',1,'IIAC金奖豆',1),(8,'2025-05-28 21:12:23.196','2025-05-28 21:12:23.197',NULL,2,'temperature','温度',1,'冰',2),(9,'2025-05-28 21:12:23.199','2025-05-28 21:12:23.216',NULL,2,'sweetness','甜度',1,'不另外加糖',3),(10,'2025-05-28 21:12:23.387','2025-05-28 21:12:23.387',NULL,3,'size','杯型',1,'大杯 16oz',0),(11,'2025-05-28 21:12:23.388','2025-05-28 21:12:23.390',NULL,3,'bean','咖啡豆',1,'IIAC金奖豆',1),(12,'2025-05-28 21:12:23.391','2025-05-28 21:12:23.392',NULL,3,'temperature','温度',1,'热',2),(13,'2025-05-28 21:12:23.392','2025-05-28 21:12:23.459',NULL,3,'sweetness','甜度',1,'不另外加糖',3),(14,'2025-05-28 21:12:23.396','2025-05-28 21:12:23.396',NULL,3,'other','奶基',1,'纯牛奶',4),(15,'2025-05-28 21:12:23.674','2025-05-28 21:12:23.675',NULL,4,'size','杯型',1,'大杯 16oz',0),(16,'2025-05-28 21:12:23.676','2025-05-28 21:12:23.676',NULL,4,'temperature','温度',1,'冰',1),(17,'2025-05-28 21:12:23.677','2025-05-28 21:12:23.699',NULL,4,'sweetness','甜度',1,'微甜',2),(18,'2025-05-28 21:12:23.679','2025-05-28 21:12:23.680',NULL,4,'flavor','口味',1,'茉莉花香',3);
/*!40000 ALTER TABLE `product_spec_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_spec_options`
--

DROP TABLE IF EXISTS `product_spec_options`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_spec_options` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `spec_id` bigint unsigned DEFAULT NULL,
  `code` varchar(50) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `is_default` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_product_spec_options_deleted_at` (`deleted_at`),
  KEY `idx_product_spec_options_spec_id` (`spec_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_spec_options`
--

LOCK TABLES `product_spec_options` WRITE;
/*!40000 ALTER TABLE `product_spec_options` DISABLE KEYS */;
INSERT INTO `product_spec_options` VALUES (1,'2025-05-28 21:12:22.746','2025-05-28 21:12:22.746',NULL,0,'0','大杯 16oz',1),(2,'2025-05-28 21:12:22.752','2025-05-28 21:12:22.752',NULL,0,'1','热',1),(3,'2025-05-28 21:12:22.754','2025-05-28 21:12:22.754',NULL,0,'2','少少甜',0),(4,'2025-05-28 21:12:22.755','2025-05-28 21:12:22.755',NULL,0,'3','微甜',0),(5,'2025-05-28 21:12:22.755','2025-05-28 21:12:22.755',NULL,0,'4','不另外加糖',1);
/*!40000 ALTER TABLE `product_spec_options` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_specs`
--

DROP TABLE IF EXISTS `product_specs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `product_specs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `sk_uid` bigint unsigned DEFAULT NULL,
  `specs_code` varchar(50) NOT NULL,
  `specs_name` varchar(100) DEFAULT NULL,
  `specs_type` bigint DEFAULT NULL,
  `is_required` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_product_specs_deleted_at` (`deleted_at`),
  KEY `idx_product_specs_sk_uid` (`sk_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_specs`
--

LOCK TABLES `product_specs` WRITE;
/*!40000 ALTER TABLE `product_specs` DISABLE KEYS */;
/*!40000 ALTER TABLE `product_specs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `goods_code` varchar(50) NOT NULL,
  `goods_name` varchar(200) NOT NULL,
  `available_specs` text,
  `last_sync_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_products_goods_code` (`goods_code`),
  KEY `idx_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES (1,'2025-05-28 21:12:22.742','2025-05-29 19:47:35.330',NULL,'2500','标准美式','{\"14\":{\"items\":[{\"code\":\"0\",\"is_default\":0,\"name\":\"标准甜\"},{\"code\":\"1\",\"is_default\":0,\"name\":\"少甜\"},{\"code\":\"2\",\"is_default\":0,\"name\":\"少少甜\"},{\"code\":\"3\",\"is_default\":0,\"name\":\"微甜\"},{\"code\":\"4\",\"is_default\":1,\"name\":\"不另外加糖\"}],\"name\":\"糖\",\"required\":true,\"type\":2},\"15\":{\"items\":[{\"code\":\"0\",\"is_default\":0,\"name\":\"双份奶\"},{\"code\":\"1\",\"is_default\":0,\"name\":\"单份奶\"},{\"code\":\"2\",\"is_default\":1,\"name\":\"无奶\"}],\"name\":\"奶\",\"required\":true,\"type\":2},\"17\":{\"items\":[{\"code\":\"0\",\"is_default\":0,\"name\":\"冰\"},{\"code\":\"1\",\"is_default\":1,\"name\":\"热\"}],\"name\":\"温度\",\"required\":true,\"type\":2},\"50\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"IIAC金奖豆\"}],\"name\":\"咖啡豆\",\"required\":true,\"type\":2},\"64\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"大杯 16oz\"}],\"name\":\"杯型\",\"required\":true,\"type\":2}}','2025-05-29 19:47:35.115'),(2,'2025-05-28 21:12:23.189','2025-05-29 19:47:35.512',NULL,'4500','燕麦拿铁','{\"17\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"冰\"},{\"code\":\"1\",\"is_default\":0,\"name\":\"热\"}],\"name\":\"温度\",\"required\":true,\"type\":2},\"18\":{\"items\":[{\"code\":\"0\",\"is_default\":0,\"name\":\"标准甜\"},{\"code\":\"1\",\"is_default\":0,\"name\":\"少甜\"},{\"code\":\"2\",\"is_default\":0,\"name\":\"少少甜\"},{\"code\":\"3\",\"is_default\":0,\"name\":\"微甜\"},{\"code\":\"4\",\"is_default\":1,\"name\":\"不另外加糖\"}],\"name\":\"糖度\",\"required\":true,\"type\":2},\"50\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"IIAC金奖豆\"}],\"name\":\"咖啡豆\",\"required\":true,\"type\":2},\"64\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"大杯 16oz\"}],\"name\":\"杯型\",\"required\":true,\"type\":2}}','2025-05-29 19:47:35.449'),(3,'2025-05-28 21:12:23.385','2025-05-29 19:47:35.789',NULL,'4805','拿铁','{\"102\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"纯牛奶\"},{\"code\":\"1\",\"is_default\":0,\"name\":\"燕麦奶\"}],\"name\":\"奶基\",\"required\":true,\"type\":2},\"17\":{\"items\":[{\"code\":\"0\",\"is_default\":0,\"name\":\"冰\"},{\"code\":\"1\",\"is_default\":1,\"name\":\"热\"}],\"name\":\"温度\",\"required\":true,\"type\":2},\"18\":{\"items\":[{\"code\":\"0\",\"is_default\":0,\"name\":\"标准甜\"},{\"code\":\"1\",\"is_default\":0,\"name\":\"少甜\"},{\"code\":\"2\",\"is_default\":0,\"name\":\"少少甜\"},{\"code\":\"3\",\"is_default\":0,\"name\":\"微甜\"},{\"code\":\"4\",\"is_default\":1,\"name\":\"不另外加糖\"}],\"name\":\"糖度\",\"required\":true,\"type\":2},\"50\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"IIAC金奖豆\"}],\"name\":\"咖啡豆\",\"required\":true,\"type\":2},\"64\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"大杯 16oz\"}],\"name\":\"杯型\",\"required\":true,\"type\":2}}','2025-05-29 19:47:35.634'),(4,'2025-05-28 21:12:23.672','2025-05-29 19:47:35.950',NULL,'4929','橙C冰茶','{\"100\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"茉莉花香\"}],\"name\":\"茶风味\",\"required\":true,\"type\":2},\"17\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"冰\"}],\"name\":\"温度\",\"required\":true,\"type\":2},\"18\":{\"items\":[{\"code\":\"0\",\"is_default\":0,\"name\":\"标准甜\"},{\"code\":\"1\",\"is_default\":0,\"name\":\"少甜\"},{\"code\":\"2\",\"is_default\":0,\"name\":\"少少甜\"},{\"code\":\"3\",\"is_default\":1,\"name\":\"微甜\"},{\"code\":\"4\",\"is_default\":0,\"name\":\"不另外加糖\"}],\"name\":\"糖度\",\"required\":true,\"type\":2},\"64\":{\"items\":[{\"code\":\"0\",\"is_default\":1,\"name\":\"大杯 16oz\"}],\"name\":\"杯型\",\"required\":true,\"type\":2}}','2025-05-29 19:47:35.921');
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `spec_input_aliases`
--

DROP TABLE IF EXISTS `spec_input_aliases`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `spec_input_aliases` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `spec_type` varchar(50) NOT NULL,
  `standard_value` varchar(50) NOT NULL,
  `alias_value` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_spec_input_aliases_deleted_at` (`deleted_at`),
  KEY `idx_spec_input_aliases_spec_type` (`spec_type`),
  KEY `idx_spec_input_aliases_alias_value` (`alias_value`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `spec_input_aliases`
--

LOCK TABLES `spec_input_aliases` WRITE;
/*!40000 ALTER TABLE `spec_input_aliases` DISABLE KEYS */;
INSERT INTO `spec_input_aliases` VALUES (1,'2025-05-29 23:26:20.872','2025-05-29 23:26:20.872',NULL,'size','大杯','大'),(2,'2025-05-29 23:26:20.880','2025-05-29 23:26:20.880',NULL,'size','大杯','16oz'),(3,'2025-05-29 23:26:20.887','2025-05-29 23:26:20.887',NULL,'size','大杯','大杯 16oz'),(4,'2025-05-29 23:26:20.895','2025-05-29 23:26:20.895',NULL,'size','中杯','中'),(5,'2025-05-29 23:26:20.904','2025-05-29 23:26:20.904',NULL,'size','中杯','12oz'),(6,'2025-05-29 23:26:20.910','2025-05-29 23:26:20.910',NULL,'size','小杯','小'),(7,'2025-05-29 23:26:20.918','2025-05-29 23:26:20.918',NULL,'size','小杯','8oz'),(8,'2025-05-29 23:26:20.928','2025-05-29 23:26:20.928',NULL,'temperature','冰','冰饮'),(9,'2025-05-29 23:26:20.936','2025-05-29 23:26:20.936',NULL,'temperature','冰','加冰'),(10,'2025-05-29 23:26:20.944','2025-05-29 23:26:20.944',NULL,'temperature','热','热饮'),(11,'2025-05-29 23:26:20.951','2025-05-29 23:26:20.951',NULL,'temperature','热','加热'),(12,'2025-05-29 23:26:20.959','2025-05-29 23:26:20.959',NULL,'sweetness','标准甜','标准'),(13,'2025-05-29 23:26:20.967','2025-05-29 23:26:20.967',NULL,'sweetness','标准甜','正常甜'),(14,'2025-05-29 23:26:20.975','2025-05-29 23:26:20.975',NULL,'sweetness','微甜','微微甜'),(15,'2025-05-29 23:26:20.984','2025-05-29 23:26:20.984',NULL,'sweetness','微甜','少糖'),(16,'2025-05-29 23:26:20.990','2025-05-29 23:26:20.990',NULL,'sweetness','少甜','半糖'),(17,'2025-05-29 23:26:21.000','2025-05-29 23:26:21.000',NULL,'sweetness','不另外加糖','无糖'),(18,'2025-05-29 23:26:21.008','2025-05-29 23:26:21.008',NULL,'sweetness','不另外加糖','不加糖');
/*!40000 ALTER TABLE `spec_input_aliases` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `system_configs`
--

DROP TABLE IF EXISTS `system_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `system_configs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `config_key` varchar(100) NOT NULL,
  `config_value` text,
  `description` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_system_configs_config_key` (`config_key`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `system_configs`
--

LOCK TABLES `system_configs` WRITE;
/*!40000 ALTER TABLE `system_configs` DISABLE KEYS */;
INSERT INTO `system_configs` VALUES (1,'2025-05-28 20:53:05.534','2025-05-29 07:30:45.678','sync_store_code','387708','商品同步和卡片验证使用的店铺代码'),(2,'2025-05-28 20:53:05.544','2025-05-29 07:30:45.686','sync_enabled','true','是否启用自动商品同步');
/*!40000 ALTER TABLE `system_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `distributor_id` bigint unsigned NOT NULL,
  `type` bigint NOT NULL,
  `amount` double NOT NULL,
  `balance_before` double NOT NULL,
  `balance_after` double NOT NULL,
  `related_id` varchar(100) DEFAULT NULL,
  `remark` varchar(500) DEFAULT NULL,
  `created_by` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_transactions_distributor_id` (`distributor_id`),
  KEY `idx_transactions_type` (`type`),
  KEY `fk_transactions_creator` (`created_by`),
  CONSTRAINT `fk_transactions_creator` FOREIGN KEY (`created_by`) REFERENCES `admins` (`id`),
  CONSTRAINT `fk_transactions_distributor` FOREIGN KEY (`distributor_id`) REFERENCES `distributors` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
INSERT INTO `transactions` VALUES (1,1,2,9.1,1000,990.9,'DD1748151760023d2d97561','订单消费：DD1748151760023d2d97561',1,'2025-05-25 13:42:40.896'),(3,3,5,50,0,50,'','余额调整测试',1,'2025-05-28 23:23:36.650'),(4,3,2,9.1,50,40.9,'DD174844582354738625164','订单消费：DD174844582354738625164',1,'2025-05-28 23:23:43.948'),(5,3,2,9.55,40.9,31.349999999999998,'DD1748526961752133ad737','订单消费：DD1748526961752133ad737',1,'2025-05-29 21:56:03.584'),(6,3,1,100,31.349999999999998,131.35,'','',1,'2025-05-29 23:36:53.463');
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `withdrawals`
--

DROP TABLE IF EXISTS `withdrawals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `withdrawals` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `distributor_id` bigint unsigned NOT NULL,
  `amount` double NOT NULL,
  `status` bigint DEFAULT '0',
  `account_info` text NOT NULL,
  `remark` varchar(500) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `processed_at` datetime(3) DEFAULT NULL,
  `processed_by` bigint unsigned DEFAULT NULL,
  `reject_reason` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_withdrawals_distributor_id` (`distributor_id`),
  KEY `idx_withdrawals_status` (`status`),
  KEY `fk_withdrawals_processor` (`processed_by`),
  CONSTRAINT `fk_withdrawals_distributor` FOREIGN KEY (`distributor_id`) REFERENCES `distributors` (`id`),
  CONSTRAINT `fk_withdrawals_processor` FOREIGN KEY (`processed_by`) REFERENCES `admins` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `withdrawals`
--

LOCK TABLES `withdrawals` WRITE;
/*!40000 ALTER TABLE `withdrawals` DISABLE KEYS */;
/*!40000 ALTER TABLE `withdrawals` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-05-30  1:39:49
