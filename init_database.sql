-- CloudEye 系统数据库初始化脚本
-- 创建日期：2025-05-21

-- 创建数据库
DROP DATABASE IF EXISTS cloud_eye;
CREATE DATABASE cloud_eye DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE cloud_eye;

-- 创建云服务商表
DROP TABLE IF EXISTS cloud_providers;
CREATE TABLE cloud_providers (
    id INT UNSIGNED AUTO_INCREMENT COMMENT '云服务商ID',
    name VARCHAR(100) NOT NULL COMMENT '云服务商名称',
    code VARCHAR(50) NOT NULL COMMENT '云服务商代码',
    description TEXT COMMENT '云服务商描述',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='云服务商信息表';

-- 创建云产品表
DROP TABLE IF EXISTS cloud_products;
CREATE TABLE cloud_products (
    id INT UNSIGNED AUTO_INCREMENT COMMENT '云产品ID',
    cloud_provider_id INT UNSIGNED NOT NULL COMMENT '关联的云服务商ID',
    name VARCHAR(100) NOT NULL COMMENT '产品名称',
    code VARCHAR(50) NOT NULL COMMENT '产品代码',
    description TEXT COMMENT '产品描述',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_provider_code (cloud_provider_id, code),
    CONSTRAINT fk_products_provider FOREIGN KEY (cloud_provider_id) REFERENCES cloud_providers (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='云产品信息表';

-- 创建安全配置基线项表
DROP TABLE IF EXISTS configuration_items;
CREATE TABLE configuration_items (
    id INT UNSIGNED AUTO_INCREMENT COMMENT '配置项ID',
    cloud_provider_id INT UNSIGNED NOT NULL COMMENT '关联的云服务商ID',
    product_id INT UNSIGNED NOT NULL COMMENT '关联的产品ID',
    name VARCHAR(200) NOT NULL COMMENT '配置项名称',
    recommended_value TEXT NOT NULL COMMENT '推荐配置值',
    risk_description TEXT COMMENT '风险说明',
    check_method TEXT COMMENT '检查方法',
    configuration_method TEXT COMMENT '配置方式',
    reference TEXT COMMENT '参考资料',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_provider_product (cloud_provider_id, product_id),
    CONSTRAINT fk_config_provider FOREIGN KEY (cloud_provider_id) REFERENCES cloud_providers (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_config_product FOREIGN KEY (product_id) REFERENCES cloud_products (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全配置基线项表';

-- 初始化云服务商数据
INSERT INTO cloud_providers (name, code, description) VALUES
    ('Amazon Web Services', 'AWS', 'Amazon Web Services (AWS) 是亚马逊（Amazon）公司旗下云计算服务平台，提供包括弹性计算、存储、数据库、机器学习等在内的一系列云服务。'),
    ('Microsoft Azure', 'AZURE', 'Microsoft Azure 是微软公司的云计算服务，为开发人员和IT专业人员构建、部署和管理应用程序提供SaaS、PaaS和IaaS等多种解决方案。'),
    ('Google Cloud Platform', 'GCP', 'Google Cloud Platform (GCP) 是由Google提供的云计算服务，包括计算、数据存储、数据分析和机器学习等一系列模块化云服务。'),
    ('阿里云', 'ALICLOUD', '阿里云是阿里巴巴集团旗下的云计算品牌，为全球企业、开发者和政府机构提供安全、可靠的计算和数据处理能力。'),
    ('腾讯云', 'TENCENTCLOUD', '腾讯云是腾讯推出的云计算品牌，提供云服务器、云存储、云数据库和大数据处理等基础云计算服务。');

-- 初始化云产品数据
INSERT INTO cloud_products (cloud_provider_id, name, code, description) VALUES
    -- AWS产品
    (1, 'Amazon Elastic Compute Cloud', 'EC2', 'Amazon EC2 是一种提供可伸缩计算容量的Web服务，让开发人员能够更轻松地进行云端计算。'),
    (1, 'Amazon Simple Storage Service', 'S3', 'Amazon S3 是一种对象存储服务，提供行业领先的可扩展性、数据可用性、安全性和性能。'),
    (1, 'Amazon Relational Database Service', 'RDS', 'Amazon RDS 让用户可在云中轻松设置、操作和扩展关系数据库。'),
    
    -- Azure产品
    (2, 'Azure Virtual Machines', 'AVM', 'Azure Virtual Machines 提供可缩放的计算资源，让用户能够灵活地运行应用程序。'),
    (2, 'Azure Blob Storage', 'BLOB', 'Azure Blob Storage 是适用于云的对象存储解决方案，用于存储大量非结构化数据。'),
    (2, 'Azure SQL Database', 'ASQL', 'Azure SQL Database 是基于最新稳定版Microsoft SQL Server数据库引擎的智能关系云数据库服务。'),
    
    -- GCP产品
    (3, 'Google Compute Engine', 'GCE', 'Google Compute Engine 提供可配置的虚拟机，在Google基础设施上运行。'),
    (3, 'Google Cloud Storage', 'GCS', 'Google Cloud Storage 是一种持久、高可用且安全的对象存储服务。'),
    (3, 'Google Cloud SQL', 'GSQL', 'Google Cloud SQL 是一种托管关系型数据库服务，用于MySQL、PostgreSQL和SQL Server。'),
    
    -- 阿里云产品
    (4, '阿里云弹性计算服务', 'ECS', '阿里云ECS是一种提供弹性可伸缩计算能力的服务，帮助用户快速构建更稳定、安全的应用。'),
    (4, '阿里云对象存储服务', 'OSS', '阿里云OSS提供海量、安全、低成本、高可靠的云存储服务，适合存储各种文件类型。'),
    (4, '阿里云关系型数据库', 'RDS', '阿里云RDS是一种稳定可靠、可弹性伸缩的在线数据库服务，提供多种数据库引擎选择。'),
    
    -- 腾讯云产品
    (5, '腾讯云服务器', 'CVM', '腾讯云CVM提供安全可靠的弹性计算服务，支持Linux、Windows等操作系统，适合承载各类应用。'),
    (5, '腾讯云对象存储', 'COS', '腾讯云COS是腾讯云提供的一种存储海量文件的分布式存储服务，具有高扩展性、低成本等优点。'),
    (5, '腾讯云数据库', 'TencentDB', '腾讯云数据库是腾讯云提供的高性能、高可靠、高安全、可弹性伸缩的数据库托管服务。');

-- 初始化安全配置基线项
INSERT INTO configuration_items (cloud_provider_id, product_id, name, recommended_value, risk_description, check_method, configuration_method, reference) VALUES
    -- AWS EC2 配置项
    (1, 1, 'EC2实例安全组入站规则限制', '仅开放必要的端口和IP范围', '不恰当的安全组规则可能导致未授权访问EC2实例上的服务。', '通过AWS控制台或CLI检查安全组规则，确保仅允许必要的入站流量。', '在AWS控制台或使用CLI修改EC2安全组规则，移除非必要的端口开放。', 'AWS安全最佳实践文档 https://docs.aws.amazon.com/security/'),
    (1, 1, 'EC2实例AMI更新状态', '使用最新的安全补丁AMI', '过时的AMI可能包含已知漏洞，增加系统被攻击的风险。', '检查AMI的创建日期和补丁级别，确保使用最新的安全补丁版本。', '定期更新EC2实例使用的AMI，或为现有实例应用安全补丁。', 'AWS AMI安全指南 https://docs.aws.amazon.com/security/ami-security/'),
    
    -- AWS S3 配置项
    (1, 2, 'S3存储桶公共访问设置', '禁用所有公共访问选项', '允许公共访问可能导致敏感数据泄露。', '使用AWS控制台或CLI检查存储桶的"阻止公共访问"设置。', '在S3存储桶配置中启用"阻止所有公共访问"选项。', 'AWS S3安全最佳实践 https://docs.aws.amazon.com/AmazonS3/latest/userguide/security-best-practices.html'),
    (1, 2, 'S3存储桶加密设置', '启用默认加密（AES-256或AWS KMS）', '未加密的数据存在被未授权访问的风险。', '检查S3存储桶的默认加密设置。', '在S3存储桶属性中启用默认加密，选择AES-256或AWS KMS。', 'AWS S3加密指南 https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucket-encryption.html'),
    
    -- AWS RDS 配置项
    (1, 3, 'RDS数据库加密设置', '启用存储加密', '未加密的数据库存储可能导致敏感信息泄露。', '检查RDS实例是否启用了存储加密。', '创建新的RDS实例时启用加密选项，或加密现有数据库的快照并从该快照恢复。', 'AWS RDS加密指南 https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Overview.Encryption.html'),
    (1, 3, 'RDS数据库公共可访问性', '禁用公共可访问性', '允许公共访问数据库增加了未授权访问的风险。', '检查RDS实例的"公共可访问性"设置。', '修改RDS实例，将"公共可访问性"设置为"否"。', 'AWS RDS安全最佳实践 https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_BestPractices.Security.html'),
    
    -- Azure VM 配置项
    (2, 4, 'Azure VM网络安全组设置', '仅允许必要的入站规则', '过于宽松的NSG规则可能导致VM被未授权访问。', '在Azure门户或使用Azure CLI检查NSG规则。', '修改NSG规则，删除非必要的入站规则，限制IP范围和端口。', 'Azure NSG安全最佳实践 https://docs.microsoft.com/azure/security/fundamentals/network-best-practices'),
    (2, 4, 'Azure VM磁盘加密', '启用Azure磁盘加密', '未加密的VM磁盘可能导致数据泄露。', '检查VM是否启用了Azure磁盘加密。', '为新VM启用磁盘加密，或对现有VM启用Azure磁盘加密。', 'Azure磁盘加密指南 https://docs.microsoft.com/azure/security/fundamentals/azure-disk-encryption-vms-vmss'),
    
    -- Azure Blob Storage 配置项
    (2, 5, '存储账户公共访问级别', '禁用公共访问', '允许公共访问可能导致数据泄露。', '检查存储账户的公共访问级别设置。', '在Azure门户中修改存储账户的"允许Blob公共访问"设置为"禁用"。', 'Azure Storage安全指南 https://docs.microsoft.com/azure/storage/blobs/security-recommendations'),
    (2, 5, '存储账户加密设置', '启用默认加密', '未加密的数据存储增加了敏感信息泄露的风险。', '检查存储账户的加密设置。', 'Azure存储账户默认启用加密，确保使用CMK（客户管理的密钥）以获得更高的安全性。', 'Azure存储加密指南 https://docs.microsoft.com/azure/storage/common/storage-service-encryption'),
    
    -- 阿里云ECS配置项
    (4, 10, 'ECS安全组规则配置', '仅开放必要的端口和授权对象', '过于宽松的安全组规则增加了被攻击的风险。', '在阿里云控制台检查安全组规则配置。', '修改安全组规则，移除不必要的入方向规则，限制端口范围和授权对象。', '阿里云安全组最佳实践 https://help.aliyun.com/document_detail/25475.html'),
    (4, 10, 'ECS实例密码复杂度', '使用高强度密码且定期更换', '弱密码容易被暴力破解，导致系统被入侵。', '检查密码策略是否符合复杂度要求。', '设置包含大小写字母、数字和特殊字符的复杂密码，定期更换。', '阿里云ECS安全最佳实践 https://help.aliyun.com/document_detail/51701.html'),
    
    -- 阿里云OSS配置项
    (4, 11, 'OSS存储桶访问控制', '使用Bucket ACL和IAM权限控制访问', '不当的访问控制可能导致数据被未授权访问。', '检查OSS Bucket的访问控制设置。', '通过OSS控制台设置合适的Bucket ACL，结合RAM权限策略控制访问。', '阿里云OSS访问控制最佳实践 https://help.aliyun.com/document_detail/31952.html');