SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for leave_type
-- ----------------------------
DROP TABLE IF EXISTS `leave_type`;
CREATE TABLE `leave_type` (
                              `id` bigint NOT NULL AUTO_INCREMENT,
                              `leave_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
                              `type_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                              `type_name_s` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                              `description` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
                              `is_reset` tinyint(1) DEFAULT '0',
                              `max_days` int DEFAULT NULL,
                              `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                              `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              PRIMARY KEY (`id`) USING BTREE,
                              UNIQUE KEY `type_name` (`type_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of leave_type
-- ----------------------------
BEGIN;
INSERT INTO `leave_type` (`id`, `leave_key`, `type_name`, `type_name_s`, `description`, `is_reset`, `max_days`, `created_at`, `updated_at`) VALUES (1, 'short_leave', 'ច្បាប់ឈប់រយៈពេលខ្លី', 'រយៈពេលខ្លី', '*មន្ត្រីនគរបាលជាតិកម្ពុជាទាំងអស់មានសិទ្ធិស្នើសុំច្បាប់ឈប់រយៈពេលខ្លីសរុបចំនួន១៥(ដប់ប្រាំ)ថ្ងៃនៃថ្ងៃធ្វើការក្នុង១(មួយ)ឆ្នាំ។', 1, 15, '2023-10-01 08:00:00', '2024-09-24 09:09:55');
INSERT INTO `leave_type` (`id`, `leave_key`, `type_name`, `type_name_s`, `description`, `is_reset`, `max_days`, `created_at`, `updated_at`) VALUES (2, 'annual_leave', 'ច្បាប់ឈប់ប្រចាំឆ្នាំ', 'ប្រចាំឆ្នាំ', '*មន្ត្រីនគរបាលជាតិដែលបានតាំងស៊ប់ក្នុងក្របខណ្ឌ មានសិទ្ធិទទួលបានច្បាប់ឈប់ប្រចាំឆ្នាំចំនួន១៥(ដប់ប្រាំ)ថ្ងៃនៃថ្ងៃធ្វើការក្នុង១(មួយ)ឆ្នាំ។', 1, 15, '2023-10-01 08:00:00', '2024-09-24 09:10:03');
INSERT INTO `leave_type` (`id`, `leave_key`, `type_name`, `type_name_s`, `description`, `is_reset`, `max_days`, `created_at`, `updated_at`) VALUES (3, 'paternity_leave', 'ច្បាប់ឈប់សម្រាកព្យាបាលជំងឺ', 'សម្រាកព្យាបាលជំងឺ', '*ច្បាប់ឈប់សម្រាកព្យាបាលជំងឺ ត្រូវបានអនុញ្ញាតឲ្យមន្ត្រីនគរបាលជាតិកម្ពុជា ពី១(មួយ)ខែទៅ៣(បី)ខែក្នុងមួយលើក។', 0, 365, '2023-10-01 08:00:00', '2024-09-24 09:10:13');
INSERT INTO `leave_type` (`id`, `leave_key`, `type_name`, `type_name_s`, `description`, `is_reset`, `max_days`, `created_at`, `updated_at`) VALUES (4, 'personal_obligation_leave', 'ច្បាប់ឈប់សម្រាកដោយមានកិច្ចការផ្ទាល់ខ្លួន', 'កិច្ចការផ្ទាល់ខ្លួន', '*ច្បាប់ឈប់សម្រាកដោយមានកិច្ចការផ្ទាល់ខ្លួន ត្រូវបានអនុញ្ញាតឲ្យមន្ត្រីនគរបាលជាតិកម្ពុជា ដើម្បីការពារផលប្រយោជន៍ផ្ទាល់ខ្លួននិងគ្រួសារ មានរយៈពេលក្នុងមួយលើកយ៉ាងតិច១(មួយ)ខែ រយៈពេលសរុបមិនឲ្យលើសពី៣(បី)ខែឡើយ ក្នុងអំឡុងពេលបម្រើការងារជាមន្ត្រីនគរបាលជាតិកម្ពុជា។', 0, 90, '2023-10-01 08:00:00', '2024-09-24 09:12:24');
INSERT INTO `leave_type` (`id`, `leave_key`, `type_name`, `type_name_s`, `description`, `is_reset`, `max_days`, `created_at`, `updated_at`) VALUES (5, 'maternity_leave', 'ច្បាប់ឈប់សម្រាកលំហែមាតុភាព', 'លំហែមាតុភាព', '*មន្ត្រីនគរបាលជាតិកម្ពុជាជាស្ត្រី​ ត្រូវបានអនុញ្ញាតច្បាប់ឈប់សម្រាកលំហែមាតុភាពរយៈពេល(បី)ខែ។', 1, 90, '2023-10-01 08:00:00', '2024-09-24 09:11:47');
COMMIT;


-- ----------------------------
-- Table structure for leave_roles
-- ----------------------------
DROP TABLE IF EXISTS `leave_roles`;
CREATE TABLE `leave_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `leave_type_id` bigint DEFAULT NULL,
  `leave_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `min_duration` int NOT NULL,
  `limit_days` int NOT NULL,
  `min_duration_show` int NOT NULL DEFAULT '1',
  `max_duration` int NOT NULL,
  `approve_level` int NOT NULL,
  `staff_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `leave_type_id` (`leave_type_id`) USING BTREE,
  CONSTRAINT `leave_roles_ibfk_1` FOREIGN KEY (`leave_type_id`) REFERENCES `leave_type` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=118 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of leave_roles
-- ----------------------------
BEGIN;
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (100, 2, 'annual_leave', 'ច្បាប់ឈប់ប្រចាំឆ្នាំ', 1, 15, 1, 1, 50, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (101, 1, 'short_leave', 'ច្បាប់ឈប់រយៈពេលខ្លី', 1, 15, 1, 1, 50, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (102, 2, 'annual_leave', 'ច្បាប់ឈប់ប្រចាំឆ្នាំ', 2, 15, 2, 2, 40, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (103, 1, 'short_leave', 'ច្បាប់ឈប់រយៈពេលខ្លី', 2, 15, 2, 2, 40, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (104, 1, 'short_leave', 'ច្បាប់ឈប់រយៈពេលខ្លី', 3, 15, 3, 4, 30, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (105, 2, 'annual_leave', 'ច្បាប់ឈប់ប្រចាំឆ្នាំ', 3, 15, 3, 4, 30, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (106, 1, 'short_leave', 'ច្បាប់ឈប់រយៈពេលខ្លី', 6, 15, 5, 11, 20, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (107, 2, 'annual_leave', 'ច្បាប់ឈប់ប្រចាំឆ្នាំ', 6, 15, 5, 11, 20, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (110, 5, 'maternity_leave', 'ច្បាប់ឈប់សម្រាកលំហែមាតុភាព', 90, 90, 90, 90, 10, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (111, 2, 'annual_leave', 'ច្បាប់ឈប់ប្រចាំឆ្នាំ', 12, 15, 12, 15, 10, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (113, 3, 'paternity_leave', 'ច្បាប់ឈប់សម្រាកព្យាបាលជំងឺ', 30, 360, 30, 90, 10, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (114, 4, 'personal_obligation_leave', 'ច្បាប់ឈប់សម្រាកដោយមានកិច្ចការផ្ទាល់ខ្លួន', 30, 90, 30, 90, 10, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (115, 1, 'short_leave', 'ច្បាប់ឈប់រយៈពេលខ្លី', 5, 15, 5, 5, 25, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (116, 2, 'annual_leave', 'ច្បាប់ឈប់ប្រចាំឆ្នាំ', 5, 15, 5, 5, 25, 'police');
INSERT INTO `leave_roles` (`id`, `leave_type_id`, `leave_type`, `name`, `min_duration`, `limit_days`, `min_duration_show`, `max_duration`, `approve_level`, `staff_type`) VALUES (117, 1, 'short_leave', 'ច្បាប់ឈប់រយៈពេលខ្លី', 12, 15, 12, 15, 10, 'police');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
