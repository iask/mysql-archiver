-- MySQL dump 10.13  Distrib 5.5.31, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: archive
-- ------------------------------------------------------
-- Server version	5.5.31-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `archive`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `archiver` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `archiver`;

--
-- Table structure for table `crontab`
--

DROP TABLE IF EXISTS `crontab`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `crontab` (
  `scheduled_id` int(10) unsigned NOT NULL,
  `scheduled_name` varchar(128) NOT NULL,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `cron` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`scheduled_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `hosts`
--

DROP TABLE IF EXISTS `hosts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hosts` (
  `host_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `hostname` varchar(255) DEFAULT NULL,
  `ip` int(10) unsigned DEFAULT NULL,
  `description` varchar(256) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `active` tinyint(4) DEFAULT '1' COMMENT '主机激活状态,是 -> 1, 否 -> 0',
  PRIMARY KEY (`host_id`),
  UNIQUE KEY `uk_hostname` (`hostname`(191))
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hosts`
--

--
-- Table structure for table `jobs`
--

DROP TABLE IF EXISTS `jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `jobs` (
  `job_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `scheduled_id` int(10) unsigned NOT NULL,
  `start_time` timestamp NOT NULL DEFAULT '2015-06-16 16:00:00',
  `running_time` int(11) DEFAULT '0',
  `end_time` datetime DEFAULT '0000-00-00 00:00:00',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1 -> initializing, 2 -> running, 3 -> completed',
  `pid` int(10) unsigned NOT NULL,
  `killed` tinyint(3) unsigned DEFAULT '0' COMMENT 'pt-archiver process status, 1 -> no, 2 -> yes',
  `target_name` varchar(256) NOT NULL DEFAULT '' COMMENT '完整的表名或是离线文件名',
  `dbhost` varchar(255) DEFAULT NULL,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `stdout` varchar(256) NOT NULL DEFAULT '' COMMENT 'stdout日志路径',
  `stderr` varchar(256) NOT NULL DEFAULT '' COMMENT 'stderr日志路径',
  `backup_id` int(11) DEFAULT '0' COMMENT '远程备份系统的任务ID',
  `backup_name` varchar(64) DEFAULT '' COMMENT '项目名',
  `backup_dir` varchar(128) DEFAULT '' COMMENT '远程备份系统目录',
  `backup_host` varchar(64) DEFAULT '' COMMENT '远程备份机',
  `backup_status` varchar(32) DEFAULT '' COMMENT '远程备份状态',
  PRIMARY KEY (`job_id`),
  KEY `idx_status_stime` (`status`,`start_time`),
  KEY `idx_sid` (`scheduled_id`)
) ENGINE=InnoDB AUTO_INCREMENT=206 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `running_jobs`
--

DROP TABLE IF EXISTS `running_jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `running_jobs` (
  `running_job_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `dbhost` varchar(255) NOT NULL,
  `scheduled_id` int(10) unsigned DEFAULT NULL,
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `pid` int(10) unsigned NOT NULL,
  `target_name` varchar(256) NOT NULL DEFAULT '' COMMENT '完整的表名或是离线文件名',
  PRIMARY KEY (`running_job_id`),
  UNIQUE KEY `i_scheduled_backup` (`scheduled_id`)
) ENGINE=InnoDB AUTO_INCREMENT=201 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `running_jobs`
--

LOCK TABLES `running_jobs` WRITE;
/*!40000 ALTER TABLE `running_jobs` DISABLE KEYS */;
/*!40000 ALTER TABLE `running_jobs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `s`
--

DROP TABLE IF EXISTS `s`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `s` (
  `scheduled_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `xboxtag` varchar(255) NOT NULL,
  `hostid` int(11) NOT NULL DEFAULT '0',
  `port` int(10) unsigned NOT NULL DEFAULT '3306',
  `db` varchar(128) NOT NULL,
  `tbl` varchar(128) NOT NULL,
  `querystr` varchar(512) NOT NULL DEFAULT '',
  `cron` varchar(128) NOT NULL DEFAULT '',
  `target_type` tinyint(4) DEFAULT '1' COMMENT '归档方式,1 -> 离线文件, 2; -> 在线表',
  `target_name` varchar(256) NOT NULL DEFAULT '' COMMENT '备用字段,默认为空',
  `charset` varchar(16) NOT NULL DEFAULT 'utf8',
  `weight_id` tinyint(4) NOT NULL DEFAULT '3',
  `active` tinyint(4) DEFAULT '0' COMMENT '任务激活状态,0 -> 未激活, 1 -> 已激活',
  `deadline` datetime NOT NULL DEFAULT '2050-01-01 00:00:00' COMMENT '该任务的有效期',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`scheduled_id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_host_port_db_table` (`xboxtag`,`port`,`db`,`tbl`),
  KEY `idx_init_job_list` (`deadline`,`active`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scheduled`
--

DROP TABLE IF EXISTS `scheduled`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scheduled` (
  `scheduled_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `xboxtag` varchar(255) NOT NULL,
  `hostid` int(11) NOT NULL DEFAULT '0',
  `port` int(10) unsigned NOT NULL DEFAULT '3306',
  `db` varchar(128) NOT NULL,
  `tbl` varchar(128) NOT NULL,
  `querystr` varchar(512) NOT NULL DEFAULT '',
  `cron` varchar(128) NOT NULL DEFAULT '',
  `target_type` tinyint(4) DEFAULT '1' COMMENT '归档方式,1 -> 离线文件, 2; -> 在线表',
  `target_name` varchar(256) NOT NULL DEFAULT '' COMMENT '备用字段,默认为空',
  `charset` varchar(16) NOT NULL DEFAULT 'utf8',
  `weight_id` tinyint(4) NOT NULL DEFAULT '3',
  `active` tinyint(4) DEFAULT '0' COMMENT '任务激活状态,0 -> 未激活, 1 -> 已激活',
  `deadline` datetime NOT NULL DEFAULT '2050-01-01 00:00:00' COMMENT '该任务的有效期',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`scheduled_id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_host_port_db_table` (`xboxtag`,`port`,`db`,`tbl`),
  KEY `idx_init_job_list` (`deadline`,`active`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `weight`
--

DROP TABLE IF EXISTS `weight`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `weight` (
  `weight_id` tinyint(3) unsigned NOT NULL AUTO_INCREMENT,
  `chunk_size` smallint(5) unsigned NOT NULL,
  `low_priority_delete` tinyint(4) NOT NULL DEFAULT '1' COMMENT '低优先级DELETE, true -> 1, false -> 0',
  `low_priority_insert` tinyint(4) NOT NULL DEFAULT '1' COMMENT '低优先级INSERT, true -> 1, false -> 0',
  PRIMARY KEY (`weight_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `weight`
--

LOCK TABLES `weight` WRITE;
/*!40000 ALTER TABLE `weight` DISABLE KEYS */;
INSERT INTO `weight` VALUES (1,200,1,1),(2,400,1,1),(3,600,1,1),(4,800,2,2),(5,1000,2,2),(6,2000,2,2);
/*!40000 ALTER TABLE `weight` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-06-03 19:07:41
