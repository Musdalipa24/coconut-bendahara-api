-- MySQL dump 10.13  Distrib 8.0.43, for Linux (x86_64)
--
-- Host: localhost    Database: bendahara
-- ------------------------------------------------------
-- Server version	8.0.43-0ubuntu0.22.04.1

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
INSERT INTO `history_transaksi` VALUES ('145f4d02-989f-497f-b105-c7fec8f83857','2025-08-16 13:26:00','sadsd','Pengeluaran',1000000),('1e302048-03a7-4f9b-af14-98dd609e03dd','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 3 dari RESKY AMALIA RUSLI (13.24.003)','Pemasukan',1),('4cda747f-a816-4039-901a-cf1152dc7fa8','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 2 dari Baharuddin Yusuf (14.25.002)','Pemasukan',7000),('734600de-4d39-4968-9880-be55c7bb4cf9','2025-08-16 04:32:00','dsds','Pemasukan',1000),('8f648400-6c7d-4e37-928e-0e0d66c32e7c','2025-08-16 03:51:00','dsads','Pemasukan',100000000),('92a7dc3d-12bd-4fd9-84ca-2e8627b4fa27','2025-08-16 05:05:00','sadsad','Pemasukan',1000),('96f6377a-cd2a-466b-9496-34369795a7e5','2025-08-08 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Nawat Sakti Al\'agasi (14.25.008)','Pemasukan',2000),('9e1d38ec-230c-4fec-9a9b-c63bc5162ce5','2025-08-15 16:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Keisya (14.25.005)','Pemasukan',7000),('a0ef5a63-c07d-4347-bab1-a47b4938e460','2025-03-01 17:00:00','Pembayaran Iuran periode 2025-02 - minggu ke 2 dari AHMAD FAISAL (13.24.001)','Pemasukan',3000),('b12f428d-2f63-4f78-939e-ee2858698570','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari AHMAD FAISAL (13.24.001)','Pemasukan',7000),('d0833bf3-6802-4cfb-b0ec-7019c545dad3','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 12 dari RESKY AMALIA RUSLI (13.24.003)','Pemasukan',7000),('e29fc613-0149-499c-acaf-9667941d990a','2025-08-16 13:25:00','fsdf','Pemasukan',1000),('eaef39e0-fd41-4854-a687-5023e5133910','2025-08-16 04:32:00','ddada','Pengeluaran',100000);
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
INSERT INTO `iuran` VALUES ('5ee7c8b2-7359-47e0-87aa-ec4db5e6dc05','2025-08',12,'2025-08-16 16:03:37'),('6c2cdc69-492d-4624-b4e6-ef33e7c1820c','2025-08',2,'2025-08-16 12:48:09'),('9dc0c25d-040e-4c0e-a476-e2f4981dc53b','2025-02',2,'2025-08-16 12:32:51'),('b08213b5-38db-4470-bcdd-584edb8bcaec','2025-08',1,'2025-08-16 05:29:19'),('d5b4c7fe-4ce1-4a19-91f5-6e0da3e34abd','2025-08',3,'2025-08-16 16:05:07');
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
INSERT INTO `laporan_keuangan` VALUES ('0f3833a6-ebbc-4732-8682-ad07be35a01e','2025-08-16 05:05:00','sadsad',1000,0,99934001,'92a7dc3d-12bd-4fd9-84ca-2e8627b4fa27'),('17923cef-61ab-4d12-9007-584376508cd7','2025-08-08 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Nawat Sakti Al\'agasi (14.25.008)',2000,0,5000,'96f6377a-cd2a-466b-9496-34369795a7e5'),('239ed9c8-a8fb-4d9f-9851-3dc5c054a36e','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 3 dari RESKY AMALIA RUSLI (13.24.003)',1,0,19001,'1e302048-03a7-4f9b-af14-98dd609e03dd'),('2e3fb5ec-01b5-49a7-9af9-0320511cd510','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 12 dari RESKY AMALIA RUSLI (13.24.003)',7000,0,26000,'d0833bf3-6802-4cfb-b0ec-7019c545dad3'),('39964e70-55df-4150-9f8b-22b179b97998','2025-08-16 04:32:00','ddada',0,100000,99933001,'eaef39e0-fd41-4854-a687-5023e5133910'),('4c22363e-674d-4cc2-ae74-3e512d996825','2025-08-16 13:25:00','fsdf',1000,0,99935001,'e29fc613-0149-499c-acaf-9667941d990a'),('513c41e3-b5a3-4028-a6e8-662aae1551fd','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 2 dari Baharuddin Yusuf (14.25.002)',7000,0,26000,'4cda747f-a816-4039-901a-cf1152dc7fa8'),('6886e6d5-23ba-426c-b6ab-472ff5254c15','2025-08-16 03:51:00','dsads',100000000,0,100033001,'8f648400-6c7d-4e37-928e-0e0d66c32e7c'),('a8ff61ab-6164-40ef-a21c-c4d8d047e598','2025-08-16 13:26:00','sadsd',0,1000000,98935001,'145f4d02-989f-497f-b105-c7fec8f83857'),('d793917c-2ba2-456f-b839-0aebd29f9e48','2025-08-16 04:32:00','dsds',1000,0,99934001,'734600de-4d39-4968-9880-be55c7bb4cf9'),('ef2c8e0d-254b-46c2-9f04-37188df53259','2025-03-01 17:00:00','Pembayaran Iuran periode 2025-02 - minggu ke 2 dari AHMAD FAISAL (13.24.001)',3000,0,3000,'a0ef5a63-c07d-4347-bab1-a47b4938e460'),('f933fae5-da59-4821-92a3-0e6528e5e5d1','2025-08-15 16:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Keisya (14.25.005)',7000,0,12000,'9e1d38ec-230c-4fec-9a9b-c63bc5162ce5'),('fafd6b90-4d3e-4b83-b391-89eaf73ef483','2025-08-15 17:00:00','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari AHMAD FAISAL (13.24.001)',7000,0,19000,'b12f428d-2f63-4f78-939e-ee2858698570');
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
INSERT INTO `member` VALUES ('0d50086e-8a1e-42c0-9002-a53ac414d2e2','14.25.011','Andika Rizky Ramadhan','anggota','2025-08-16 19:43:10','2025-08-16 19:43:10'),('10a26ee5-70eb-415f-bd58-ee6fb9361bf6','14.25.002','Baharuddin Yusuf','anggota','2025-08-16 19:35:37','2025-08-16 19:35:37'),('2d41d697-cbdb-465f-a1db-2413cdada748','14.25.005','Keisya','anggota','2025-08-16 13:28:46','2025-08-16 13:28:46'),('4b3c7863-27b5-42bc-8167-30f2e269fea8','14.25.008','Nawat Sakti Al\'agasi','anggota','2025-08-16 19:41:56','2025-08-16 19:41:56'),('6142716b-6c1d-4f9e-b71c-34d64fa9f83f','14.25.013','Naufal Asila','anggota','2025-08-16 19:44:56','2025-08-16 19:44:56'),('6fa7fcc4-ffba-46e9-b977-2dda5418a352','14.25.006','Bayyin Ramadhan','anggota','2025-08-16 19:41:04','2025-08-16 19:41:04'),('9cd5aeff-79d2-11f0-a92b-482ae3455d6d','13.24.001','AHMAD FAISAL','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5d011-79d2-11f0-a92b-482ae3455d6d','13.24.003','RESKY AMALIA RUSLI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5d6e8-79d2-11f0-a92b-482ae3455d6d','13.24.004','SALSABILA PUTRI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5dda2-79d2-11f0-a92b-482ae3455d6d','13.24.005','SYAHRUL RAMADHAN','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5e29b-79d2-11f0-a92b-482ae3455d6d','13.24.006','ANDI CITRA AYU LESTARI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5f36f-79d2-11f0-a92b-482ae3455d6d','13.24.007','PARWATI','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd5fb95-79d2-11f0-a92b-482ae3455d6d','13.24.008','MUHAMMAD SYARIF','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd60b85-79d2-11f0-a92b-482ae3455d6d','13.24.009','YUSUF MARCELINO ISHAK','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd63342-79d2-11f0-a92b-482ae3455d6d','13.24.011','MUH. FIKRI HAEKAL','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd63b78-79d2-11f0-a92b-482ae3455d6d','13.24.014','MUSDALIPA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd63e42-79d2-11f0-a92b-482ae3455d6d','13.24.015','MUHAMMMAD AKSAN','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd64220-79d2-11f0-a92b-482ae3455d6d','13.24.016','MUSTIKA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd645f0-79d2-11f0-a92b-482ae3455d6d','13.24.017','WINDU YOGA NUGRAHA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('9cd64d72-79d2-11f0-a92b-482ae3455d6d','13.24.018','AMELIA','bph','2025-08-15 20:23:12','2025-08-15 20:23:12'),('ab16e3cf-cfeb-47e9-b31f-9c7b17a27a56','14.25.003','Ahmad Fajrul I.','anggota','2025-08-16 19:39:58','2025-08-16 19:39:58'),('b180e799-e8ec-4daf-acd0-93f3d958952a','14.25.009','Saudah Al','anggota','2025-08-16 19:42:14','2025-08-16 19:42:14'),('b9a6b1cb-0d58-45ff-80a4-8bd0fd7558c9','14.25.010','Rizky Akbar','anggota','2025-08-16 19:42:38','2025-08-16 19:42:38'),('c1f4c52a-cb68-49ad-9e93-961065d46dc4','14.25.012','Destyna Auliany','anggota','2025-08-16 19:44:05','2025-08-16 19:44:05'),('c277b8d1-6392-4be8-ad9d-5f2f96723ba9','14.25.001','Indal Awalaikal','anggota','2025-08-16 13:27:23','2025-08-16 13:27:23'),('cc7f4c16-62a6-4efe-8f6c-aae24d046c26','14.25.007','Nurhalisa','anggota','2025-08-16 19:41:27','2025-08-16 19:41:27'),('e88e471b-6dcb-49ff-8322-4bf767d8cdc8','14.25.004','Nurhasanah','anggota','2025-08-16 19:40:24','2025-08-16 19:40:24');
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
INSERT INTO `pemasukan` VALUES ('2e210abd-5e88-48a9-ae39-e5b9e3b83461','2025-08-16 03:51:00','Dana Desa','dsads',100000000,'2025-08-16-11-51-f733bce6-35f6-458e-805d-553a71704dbc.png','8f648400-6c7d-4e37-928e-0e0d66c32e7c'),('4fe3616d-3218-4213-9069-01721816629d','2025-08-16 05:05:00','Pajak','sadsad',1000,'','92a7dc3d-12bd-4fd9-84ca-2e8627b4fa27'),('524c52a4-4757-4bf6-a8a2-bd2997fc964c','2025-08-15 17:00:00','Iuran','Pembayaran Iuran periode 2025-08 - minggu ke 12 dari RESKY AMALIA RUSLI (13.24.003)',7000,'no data','d0833bf3-6802-4cfb-b0ec-7019c545dad3'),('568040f5-3559-419a-8294-65f936e1c7e6','2025-08-15 16:00:00','Iuran','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Keisya (14.25.005)',7000,'no data','9e1d38ec-230c-4fec-9a9b-c63bc5162ce5'),('898d4235-8daa-4166-82fa-ddd7fd5f73f7','2025-08-15 17:00:00','Iuran','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari AHMAD FAISAL (13.24.001)',7000,'no data','b12f428d-2f63-4f78-939e-ee2858698570'),('8dec9d2e-97d3-4975-a85e-99430b6410d7','2025-08-15 17:00:00','Iuran','Pembayaran Iuran periode 2025-08 - minggu ke 2 dari Baharuddin Yusuf (14.25.002)',7000,'no data','4cda747f-a816-4039-901a-cf1152dc7fa8'),('a683e8f4-ce14-4979-9474-11d3a75e7fe0','2025-03-01 17:00:00','Iuran','Pembayaran Iuran periode 2025-02 - minggu ke 2 dari AHMAD FAISAL (13.24.001)',3000,'no data','a0ef5a63-c07d-4347-bab1-a47b4938e460'),('cd55a05b-f6a1-435c-95e8-bb351d64703b','2025-08-08 17:00:00','Iuran','Pembayaran Iuran periode 2025-08 - minggu ke 1 dari Nawat Sakti Al\'agasi (14.25.008)',2000,'no data','96f6377a-cd2a-466b-9496-34369795a7e5'),('ce5188b6-7e5b-4ff1-94c2-8391ae788117','2025-08-16 04:32:00','Retribusi','dsds',1000,'2025-08-16-12-32-2ca19a55-4b77-414d-9692-b6b549e6f91f.png','734600de-4d39-4968-9880-be55c7bb4cf9'),('e4fc4897-c81b-462a-a294-2a6621069b88','2025-08-15 17:00:00','Iuran','Pembayaran Iuran periode 2025-08 - minggu ke 3 dari RESKY AMALIA RUSLI (13.24.003)',1,'no data','1e302048-03a7-4f9b-af14-98dd609e03dd'),('eee7efbf-5640-4f65-aec4-c5c00f5af7c9','2025-08-16 13:25:00','sfsd','fsdf',1000,'2025-08-16-20-25-ad8256cf-30e7-4367-bda3-3d79d4d7394d.png','e29fc613-0149-499c-acaf-9667941d990a');
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
INSERT INTO `pembayaran_iuran` VALUES ('09f2e961-c901-4693-8f31-7feb45e93fec','9cd5d011-79d2-11f0-a92b-482ae3455d6d','d5b4c7fe-4ce1-4a19-91f5-6e0da3e34abd','e4fc4897-c81b-462a-a294-2a6621069b88','belum',1,'2025-08-16 00:00:00','2025-08-16 23:05:07','2025-08-16 23:05:07'),('450fad11-ae71-4c46-a812-092fcd8701c6','2d41d697-cbdb-465f-a1db-2413cdada748','b08213b5-38db-4470-bcdd-584edb8bcaec','568040f5-3559-419a-8294-65f936e1c7e6','lunas',7000,'2025-08-16 00:00:00','2025-08-16 13:29:19','2025-08-16 13:29:19'),('5e26f80b-1a5b-44a0-81bf-98c349fd96e3','9cd5aeff-79d2-11f0-a92b-482ae3455d6d','9dc0c25d-040e-4c0e-a476-e2f4981dc53b','a683e8f4-ce14-4979-9474-11d3a75e7fe0','belum',3000,'2025-03-02 00:00:00','2025-08-16 19:32:51','2025-08-16 19:32:51'),('68a794b6-4982-4db0-b14f-2b70a38a2671','9cd5d011-79d2-11f0-a92b-482ae3455d6d','5ee7c8b2-7359-47e0-87aa-ec4db5e6dc05','524c52a4-4757-4bf6-a8a2-bd2997fc964c','lunas',7000,'2025-08-16 00:00:00','2025-08-16 23:03:38','2025-08-16 23:03:38'),('7bdc7477-1071-4ab4-91c2-b3453923cd7c','10a26ee5-70eb-415f-bd58-ee6fb9361bf6','6c2cdc69-492d-4624-b4e6-ef33e7c1820c','8dec9d2e-97d3-4975-a85e-99430b6410d7','lunas',7000,'2025-08-16 00:00:00','2025-08-16 19:48:09','2025-08-16 19:48:09'),('99a440a1-1a0e-4a9b-ab01-88c943fb99d9','9cd5aeff-79d2-11f0-a92b-482ae3455d6d','b08213b5-38db-4470-bcdd-584edb8bcaec','898d4235-8daa-4166-82fa-ddd7fd5f73f7','lunas',7000,'2025-08-16 00:00:00','2025-08-16 19:26:24','2025-08-16 19:26:24'),('ef8d9507-2570-4693-8fc4-bb705cfdd024','4b3c7863-27b5-42bc-8167-30f2e269fea8','b08213b5-38db-4470-bcdd-584edb8bcaec','cd55a05b-f6a1-435c-95e8-bb351d64703b','belum',2000,'2025-08-09 00:00:00','2025-08-16 19:49:41','2025-08-16 19:49:41');
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
INSERT INTO `pengeluaran` VALUES ('0432cd65-ad0c-4e90-95d3-261a7a05ecf4','2025-08-16 13:26:00','2025-08-16-20-26-72635fcd-bdb0-4e3e-9436-e1bc71f3d51f.jpeg',1000000,'sadsd','145f4d02-989f-497f-b105-c7fec8f83857'),('d9ba600d-f13a-4bf8-a6f1-a91270a38692','2025-08-16 04:32:00','2025-08-16-12-32-d08b16ce-2988-4c89-be71-3d86c8d34b86.jpeg',100000,'ddada','eaef39e0-fd41-4854-a687-5023e5133910');
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

-- Dump completed on 2025-08-26 22:28:43
