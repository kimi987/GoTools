package sqldb

const initSql = `

-- ----------------------------
--  Table structure for ` + "`" + `guild` + "`" + `
-- ----------------------------
DROP TABLE IF EXISTS ` + "`" + `guild` + "`" + `;
CREATE TABLE ` + "`" + `guild` + "`" + ` (
  ` + "`" + `id` + "`" + ` bigint(20) NOT NULL,
  ` + "`" + `data` + "`" + ` mediumblob NOT NULL,
  PRIMARY KEY (` + "`" + `id` + "`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for ` + "`" + `hero` + "`" + `
-- ----------------------------
DROP TABLE IF EXISTS ` + "`" + `hero` + "`" + `;
CREATE TABLE ` + "`" + `hero` + "`" + ` (
  ` + "`" + `id` + "`" + ` bigint(20) NOT NULL,
  ` + "`" + `name` + "`" + ` varchar(50) DEFAULT NULL,
  ` + "`" + `hero_data` + "`" + ` mediumblob,
  ` + "`" + `base_region` + "`" + ` bigint(20) NOT NULL,
  PRIMARY KEY (` + "`" + `id` + "`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for ` + "`" + `kv` + "`" + `
-- ----------------------------
DROP TABLE IF EXISTS ` + "`" + `kv` + "`" + `;
CREATE TABLE ` + "`" + `kv` + "`" + ` (
  ` + "`" + `k` + "`" + ` varchar(20) NOT NULL,
  ` + "`" + `v` + "`" + ` mediumblob NOT NULL,
  PRIMARY KEY (` + "`" + `k` + "`" + `)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for ` + "`" + `mail` + "`" + `
-- ----------------------------
DROP TABLE IF EXISTS ` + "`" + `mail` + "`" + `;
CREATE TABLE ` + "`" + `mail` + "`" + ` (
  ` + "`" + `id` + "`" + ` bigint(20) unsigned NOT NULL,
  ` + "`" + `receiver` + "`" + ` bigint(20) NOT NULL,
  ` + "`" + `data` + "`" + ` mediumblob NOT NULL,
  ` + "`" + `keep` + "`" + ` tinyint(1) NOT NULL,
  ` + "`" + `readed` + "`" + ` tinyint(1) NOT NULL,
  ` + "`" + `has_report` + "`" + ` tinyint(1) NOT NULL,
  ` + "`" + `has_prize` + "`" + ` tinyint(1) NOT NULL,
  ` + "`" + `collected` + "`" + ` tinyint(1) NOT NULL,
  ` + "`" + `time` + "`" + ` bigint(20) NOT NULL,
  PRIMARY KEY (` + "`" + `id` + "`" + `),
  KEY ` + "`" + `receiver` + "`" + ` (` + "`" + `receiver` + "`" + `) USING HASH,
  KEY ` + "`" + `time` + "`" + ` (` + "`" + `time` + "`" + `) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
--  Table structure for ` + "`" + `dbversion` + "`" + `
-- ----------------------------
DROP TABLE IF EXISTS ` + "`" + `dbversion` + "`" + `;
CREATE TABLE ` + "`" + `dbversion` + "`" + ` (
  ` + "`" + `v` + "`" + ` int(11) unsigned NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into dbversion(v) values(0);
`
