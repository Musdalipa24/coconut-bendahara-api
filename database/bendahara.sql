-- MySQL dump 10.13  Distrib 8.0.42, for Win64 (x86_64)
--
-- Host: localhost    Database: bendahara
-- ------------------------------------------------------
-- Server version	8.0.42

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
-- Table structure for table `admin`
--

DROP TABLE IF EXISTS `admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admin` (
  `id` varchar(65) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `nik` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `username` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `role` enum('superAdmin','bendahara','guest') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admin`
--

LOCK TABLES `admin` WRITE;
/*!40000 ALTER TABLE `admin` DISABLE KEYS */;
INSERT INTO `admin` VALUES ('4433c69f-2003-42a7-9676-ea9b9dbc9f33','123','admin','$2a$10$IJxPl8c3xL4guTA/XiTFJ.Rhn5Mqi4bz.ZkcLv5c3y4JUJ8CsvBtm','superAdmin');
/*!40000 ALTER TABLE `admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `history_transaksi`
--

DROP TABLE IF EXISTS `history_transaksi`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `history_transaksi` (
  `id_transaksi` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `keterangan` varchar(255) DEFAULT NULL,
  `jenis_transaksi` enum('Pemasukan','Pengeluaran') NOT NULL,
  `nominal` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_transaksi`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `history_transaksi`
--

LOCK TABLES `history_transaksi` WRITE;
/*!40000 ALTER TABLE `history_transaksi` DISABLE KEYS */;
INSERT INTO `history_transaksi` VALUES ('734600de-4d39-4968-9880-be55c7bb4cf9','2025-08-16 04:32:00','dsds','Pemasukan',1000),('8f648400-6c7d-4e37-928e-0e0d66c32e7c','2025-08-16 03:51:00','dsads','Pemasukan',100000000),('92a7dc3d-12bd-4fd9-84ca-2e8627b4fa27','2025-08-16 05:05:00','sadsad','Pemasukan',1000),('9e1d38ec-230c-4fec-9a9b-c63bc5162ce5','2025-08-15 16:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Keisya (14.25.005)','Pemasukan',7000),('eaef39e0-fd41-4854-a687-5023e5133910','2025-08-16 04:32:00','ddada','Pengeluaran',100000);
/*!40000 ALTER TABLE `history_transaksi` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `iuran`
--

DROP TABLE IF EXISTS `iuran`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `iuran` (
  `id_iuran` varchar(36) NOT NULL,
  `periode` varchar(36) NOT NULL,
  `minggu_ke` tinyint unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_iuran`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `iuran`
--

LOCK TABLES `iuran` WRITE;
/*!40000 ALTER TABLE `iuran` DISABLE KEYS */;
INSERT INTO `iuran` VALUES ('b08213b5-38db-4470-bcdd-584edb8bcaec','2025-08',1,'2025-08-16 05:29:19');
/*!40000 ALTER TABLE `iuran` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `laporan_keuangan`
--

DROP TABLE IF EXISTS `laporan_keuangan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `laporan_keuangan` (
  `id_laporan` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `keterangan` varchar(255) DEFAULT NULL,
  `pemasukan` bigint NOT NULL,
  `pengeluaran` bigint NOT NULL,
  `saldo` bigint NOT NULL,
  `id_transaksi` varchar(36) NOT NULL,
  PRIMARY KEY (`id_laporan`),
  KEY `id_transaksi` (`id_transaksi`),
  CONSTRAINT `laporan_keuangan_ibfk_1` FOREIGN KEY (`id_transaksi`) REFERENCES `history_transaksi` (`id_transaksi`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `laporan_keuangan`
--

LOCK TABLES `laporan_keuangan` WRITE;
/*!40000 ALTER TABLE `laporan_keuangan` DISABLE KEYS */;
INSERT INTO `laporan_keuangan` VALUES ('0f3833a6-ebbc-4732-8682-ad07be35a01e','2025-08-16 05:05:00','sadsad',1000,0,99908000,'92a7dc3d-12bd-4fd9-84ca-2e8627b4fa27'),('39964e70-55df-4150-9f8b-22b179b97998','2025-08-16 04:32:00','ddada',0,100000,99907000,'eaef39e0-fd41-4854-a687-5023e5133910'),('6886e6d5-23ba-426c-b6ab-472ff5254c15','2025-08-16 03:51:00','dsads',100000000,0,100007000,'8f648400-6c7d-4e37-928e-0e0d66c32e7c'),('d793917c-2ba2-456f-b839-0aebd29f9e48','2025-08-16 04:32:00','dsds',1000,0,99908000,'734600de-4d39-4968-9880-be55c7bb4cf9'),('f933fae5-da59-4821-92a3-0e6528e5e5d1','2025-08-15 16:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Keisya (14.25.005)',7000,0,7000,'9e1d38ec-230c-4fec-9a9b-c63bc5162ce5');
/*!40000 ALTER TABLE `laporan_keuangan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `member`
--

DROP TABLE IF EXISTS `member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `member` (
  `id_member` varchar(36) NOT NULL,
  `nra` varchar(10) DEFAULT NULL,
  `nama` varchar(100) NOT NULL,
  `status` enum('bph','anggota') DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_member`),
  UNIQUE KEY `nra` (`nra`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `member`
--

LOCK TABLES `member` WRITE;
/*!40000 ALTER TABLE `member` DISABLE KEYS */;
INSERT INTO `member` VALUES ('2d41d697-cbdb-465f-a1db-2413cdada748','14.25.005','Keisya','anggota','2025-08-16 13:28:46','2025-08-16 13:28:46'),('9cd5aeff-79d2-11f0-a92b-482ae3455d6d','13.24.001','AHMAD FAISAL','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5d011-79d2-11f0-a92b-482ae3455d6d','13.24.003','RESKY AMALIA RUSLI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5d6e8-79d2-11f0-a92b-482ae3455d6d','13.24.004','SALSABILA PUTRI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5dda2-79d2-11f0-a92b-482ae3455d6d','13.24.005','SYAHRUL RAMADHAN','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5e29b-79d2-11f0-a92b-482ae3455d6d','13.24.006','ANDI CITRA AYU LESTARI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5f36f-79d2-11f0-a92b-482ae3455d6d','13.24.007','PARWATI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5fb95-79d2-11f0-a92b-482ae3455d6d','13.24.008','MUHAMMAD SYARIF','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd60b85-79d2-11f0-a92b-482ae3455d6d','13.24.009','YUSUF MARCELINO ISHAK','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd63342-79d2-11f0-a92b-482ae3455d6d','13.24.011','MUH. FIKRI HAEKAL','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd63b78-79d2-11f0-a92b-482ae3455d6d','13.24.014','MUSDALIPA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd63e42-79d2-11f0-a92b-482ae3455d6d','13.24.015','MUHAMMMAD AKSAN','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd64220-79d2-11f0-a92b-482ae3455d6d','13.24.016','MUSTIKA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd645f0-79d2-11f0-a92b-482ae3455d6d','13.24.017','WINDU YOGA NUGRAHA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd64d72-79d2-11f0-a92b-482ae3455d6d','13.24.018','AMELIA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('c277b8d1-6392-4be8-ad9d-5f2f96723ba9','14.25.001','Indal Awalaikal','anggota','2025-08-16 13:27:23','2025-08-16 13:27:23');
/*!40000 ALTER TABLE `member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pemasukan`
--

DROP TABLE IF EXISTS `pemasukan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pemasukan` (
  `id_pemasukan` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `kategori` varchar(255) NOT NULL,
  `keterangan` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `nominal` bigint NOT NULL,
  `nota` varchar(255) DEFAULT 'no data',
  `id_transaksi` varchar(36) NOT NULL,
  PRIMARY KEY (`id_pemasukan`),
  KEY `id_transaksi` (`id_transaksi`),
  CONSTRAINT `pemasukan_ibfk_1` FOREIGN KEY (`id_transaksi`) REFERENCES `history_transaksi` (`id_transaksi`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pemasukan`
--

LOCK TABLES `pemasukan` WRITE;
/*!40000 ALTER TABLE `pemasukan` DISABLE KEYS */;
INSERT INTO `pemasukan` VALUES ('2e210abd-5e88-48a9-ae39-e5b9e3b83461','2025-08-16 03:51:00','Dana Desa','dsads',100000000,'2025-08-16-11-51-f733bce6-35f6-458e-805d-553a71704dbc.png','8f648400-6c7d-4e37-928e-0e0d66c32e7c'),('4fe3616d-3218-4213-9069-01721816629d','2025-08-16 05:05:00','Pajak','sadsad',1000,'','92a7dc3d-12bd-4fd9-84ca-2e8627b4fa27'),('568040f5-3559-419a-8294-65f936e1c7e6','2025-08-15 16:00:00','Iuran','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Keisya (14.25.005)',7000,'no data','9e1d38ec-230c-4fec-9a9b-c63bc5162ce5'),('ce5188b6-7e5b-4ff1-94c2-8391ae788117','2025-08-16 04:32:00','Retribusi','dsds',1000,'2025-08-16-12-32-2ca19a55-4b77-414d-9692-b6b549e6f91f.png','734600de-4d39-4968-9880-be55c7bb4cf9');
/*!40000 ALTER TABLE `pemasukan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pembayaran_iuran`
--

DROP TABLE IF EXISTS `pembayaran_iuran`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pembayaran_iuran` (
  `id_pembayaran` varchar(36) NOT NULL,
  `id_member` varchar(36) NOT NULL,
  `id_iuran` varchar(36) NOT NULL,
  `id_pemasukan` varchar(36) NOT NULL,
  `status` enum('lunas','belum') DEFAULT 'belum',
  `jumlah_bayar` int unsigned NOT NULL,
  `tanggal_bayar` datetime DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_pembayaran`),
  KEY `id_member` (`id_member`),
  KEY `id_iuran` (`id_iuran`),
  CONSTRAINT `pembayaran_iuran_ibfk_1` FOREIGN KEY (`id_member`) REFERENCES `member` (`id_member`) ON DELETE CASCADE,
  CONSTRAINT `pembayaran_iuran_ibfk_2` FOREIGN KEY (`id_iuran`) REFERENCES `iuran` (`id_iuran`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pembayaran_iuran`
--

LOCK TABLES `pembayaran_iuran` WRITE;
/*!40000 ALTER TABLE `pembayaran_iuran` DISABLE KEYS */;
INSERT INTO `pembayaran_iuran` VALUES ('450fad11-ae71-4c46-a812-092fcd8701c6','2d41d697-cbdb-465f-a1db-2413cdada748','b08213b5-38db-4470-bcdd-584edb8bcaec','568040f5-3559-419a-8294-65f936e1c7e6','lunas',7000,'2025-08-16 00:00:00','2025-08-16 13:29:19','2025-08-16 13:29:19');
/*!40000 ALTER TABLE `pembayaran_iuran` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pengeluaran`
--

DROP TABLE IF EXISTS `pengeluaran`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pengeluaran` (
  `id_pengeluaran` varchar(36) NOT NULL,
  `tanggal` timestamp NOT NULL,
  `nota` varchar(255) NOT NULL,
  `nominal` bigint NOT NULL,
  `keterangan` varchar(255) DEFAULT NULL,
  `id_transaksi` varchar(36) NOT NULL,
  PRIMARY KEY (`id_pengeluaran`),
  KEY `id_transaksi` (`id_transaksi`),
  CONSTRAINT `pengeluaran_ibfk_1` FOREIGN KEY (`id_transaksi`) REFERENCES `history_transaksi` (`id_transaksi`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pengeluaran`
--

LOCK TABLES `pengeluaran` WRITE;
/*!40000 ALTER TABLE `pengeluaran` DISABLE KEYS */;
INSERT INTO `pengeluaran` VALUES ('d9ba600d-f13a-4bf8-a6f1-a91270a38692','2025-08-16 04:32:00','2025-08-16-12-32-d08b16ce-2988-4c89-be71-3d86c8d34b86.jpeg',100000,'ddada','eaef39e0-fd41-4854-a687-5023e5133910');
/*!40000 ALTER TABLE `pengeluaran` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-08-16 18:50:23
