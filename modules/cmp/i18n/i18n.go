// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func InitI18N() {
	// ops record
	message.SetString(language.SimplifiedChinese, "addAliECSEdgeCluster", "添加阿里云集群")
	message.SetString(language.SimplifiedChinese, "addAliACKEdgeCluster", "添加阿里云容器服务集群") // TODO remove
	message.SetString(language.SimplifiedChinese, "addAliCSEdgeCluster", "添加阿里云容器服务标准专有集群")
	message.SetString(language.SimplifiedChinese, "addAliCSManagedEdgeCluster", "添加阿里云容器服务标准托管集群")
	message.SetString(language.SimplifiedChinese, "upgradeEdgeCluster", "升级边缘集群")
	message.SetString(language.SimplifiedChinese, "offlineEdgeCluster", "集群下线")
	message.SetString(language.SimplifiedChinese, "importKubernetesCluster", "导入Kubernetes集群")
	message.SetString(language.SimplifiedChinese, "buyNodes", "购买云资源")
	message.SetString(language.SimplifiedChinese, "addNodes", "添加机器")
	message.SetString(language.SimplifiedChinese, "addEssNodes", "弹性扩容")
	message.SetString(language.SimplifiedChinese, "addAliNodes", "添加阿里云机器")
	message.SetString(language.SimplifiedChinese, "rmNodes", "下线机器")
	message.SetString(language.SimplifiedChinese, "deleteNodes", "删除机器")
	message.SetString(language.SimplifiedChinese, "deleteEssNodes", "弹性缩容")
	message.SetString(language.SimplifiedChinese, "deleteEssNodesCronJob", "弹性缩容定时任务")
	message.SetString(language.SimplifiedChinese, "setLabels", "设置标签")

	message.SetString(language.SimplifiedChinese, "createAliCloudMysql", "添加阿里云数据库实例")
	message.SetString(language.SimplifiedChinese, "createAliCloudMysqlDB", "添加阿里云数据库DB")
	message.SetString(language.SimplifiedChinese, "createAliCloudRedis", "添加阿里云Redis")
	message.SetString(language.SimplifiedChinese, "createAliCloudOss", "添加阿里云存储桶OSS")
	message.SetString(language.SimplifiedChinese, "createAliCloudOns", "添加阿里云消息队列MQ")
	message.SetString(language.SimplifiedChinese, "createAliCloudOnsTopic", "添加阿里云消息队列Topic")

	message.SetString(language.SimplifiedChinese, "DCOS-compatible label", "DCOS 兼容性的 label")
	message.SetString(language.SimplifiedChinese, "cluster name", "集群标识")
	message.SetString(language.SimplifiedChinese, "cluster type", "集群类型")
	message.SetString(language.SimplifiedChinese, "cluster version", "集群版本")
	message.SetString(language.SimplifiedChinese, "root domain", "泛域名")
	message.SetString(language.SimplifiedChinese, "edge cluster", "边缘集群")
	message.SetString(language.SimplifiedChinese, "master num", "master 数")
	message.SetString(language.SimplifiedChinese, "lb num", "lb 数")
	message.SetString(language.SimplifiedChinese, "basic", "基本信息")
	message.SetString(language.SimplifiedChinese, "init", "初始化")
	message.SetString(language.SimplifiedChinese, "plan", "执行计划")
	message.SetString(language.SimplifiedChinese, "diceInstall", "安装dice")
	message.SetString(language.SimplifiedChinese, "completed", "完成")
	message.SetString(language.SimplifiedChinese, "https enabled", "开启https")
	message.SetString(language.SimplifiedChinese, "cluster status", "集群状态")
	message.SetString(language.SimplifiedChinese, "init job cluster name", "初始化任务集群")
	message.SetString(language.SimplifiedChinese, "manage type", "管理方式")
	message.SetString(language.SimplifiedChinese, "node count", "节点数")
	message.SetString(language.SimplifiedChinese, "cluster display name", "集群名称")
	message.SetString(language.SimplifiedChinese, "cluster display name", "集群名称")
	message.SetString(language.SimplifiedChinese, "cluster init container id", "集群初始化任务容器ID")
	// cloud resource type
	message.SetString(language.SimplifiedChinese, "Compute", "计算")
	message.SetString(language.SimplifiedChinese, "Network", "网络")
	message.SetString(language.SimplifiedChinese, "Storage", "存储")
	message.SetString(language.SimplifiedChinese, "Addon", "中间件")
	// cloud resource
	message.SetString(language.SimplifiedChinese, "ECS", "云服务器 ECS")
	message.SetString(language.SimplifiedChinese, "Container Service for Kubernetes", "容器服务 Kubernetes版")
	message.SetString(language.SimplifiedChinese, "VPC", "专有网络 VPC")
	message.SetString(language.SimplifiedChinese, "SLB", "负载均衡 SLB")
	message.SetString(language.SimplifiedChinese, "EIP", "弹性公网 IP")
	message.SetString(language.SimplifiedChinese, "NAT Gateway", "NAT 网关")
	message.SetString(language.SimplifiedChinese, "NAS Storage", "NAS 文件存储")
	message.SetString(language.SimplifiedChinese, "RDS", "云数据库 RDS")
	message.SetString(language.SimplifiedChinese, "ElasticSearch", "ElasticSearch")
	// disk spec
	message.SetString(language.SimplifiedChinese, "System Disk: cloud_ssd, 200G", "系统盘: 200G SSD云盘")
	message.SetString(language.SimplifiedChinese, "System Disk: cloud_ssd, 40G", "系统盘: 40G SSD云盘")
	message.SetString(language.SimplifiedChinese, "Data Disk: cloud_ssd, 200G", "数据盘: 200G SSD云盘")
	// nat spec
	message.SetString(language.SimplifiedChinese, "Small", "小型")
	message.SetString(language.SimplifiedChinese, "Inbound Bandwidth: 10M", "入口带宽: 10M")
	message.SetString(language.SimplifiedChinese, "Outbound Bandwidth: 100M", "出口带宽: 100M")
	// others
	message.SetString(language.SimplifiedChinese, "others", "其他")
	message.SetString(language.SimplifiedChinese, "service", "服务")
	message.SetString(language.SimplifiedChinese, "env", "环境")
	message.SetString(language.SimplifiedChinese, "platform", "平台")
	message.SetString(language.SimplifiedChinese, "job", "任务")
	// cloud resource
	message.SetString(language.SimplifiedChinese, "Description", "说明")
	message.SetString(language.SimplifiedChinese, "Spec", "规格")
	message.SetString(language.SimplifiedChinese, "Status", "状态")
	message.SetString(language.SimplifiedChinese, "ChargeType", "付费方式")
	message.SetString(language.SimplifiedChinese, "Version", "版本")
	message.SetString(language.SimplifiedChinese, "NodeSpec", "节点规格")
	message.SetString(language.SimplifiedChinese, "NodeCount", "节点数")
	message.SetString(language.SimplifiedChinese, "UsedSize", "使用量")
	message.SetString(language.SimplifiedChinese, "Region", "区域(Region)")
	message.SetString(language.SimplifiedChinese, "Address", "地址")
	message.SetString(language.SimplifiedChinese, "IpAddress", "ip地址")
	message.SetString(language.SimplifiedChinese, "Bandwidth", "带宽")
	message.SetString(language.SimplifiedChinese, "Name", "名称")
	message.SetString(language.SimplifiedChinese, "NatGatewayCount", "Nat网关数量")
	message.SetString(language.SimplifiedChinese, "CurrentVersion", "当前版本")
	message.SetString(language.SimplifiedChinese, "State", "状态")
	message.SetString(language.SimplifiedChinese, "Size", "规格")
	message.SetString(language.SimplifiedChinese, "ClusterType", "集群类型")
	message.SetString(language.SimplifiedChinese, "PrePaid", "包年包月")
	message.SetString(language.SimplifiedChinese, "PostPaid", "按量付费")
	// cloud resource status
	message.SetString(language.SimplifiedChinese, "Associating", "绑定中")
	message.SetString(language.SimplifiedChinese, "Unassociating", "解绑中")
	message.SetString(language.SimplifiedChinese, "InUse", "已分配")
	message.SetString(language.SimplifiedChinese, "Available", "可用")
	message.SetString(language.SimplifiedChinese, "Pending", "等待")
	message.SetString(language.SimplifiedChinese, "Running", "运行中")
	message.SetString(language.SimplifiedChinese, "Stopping", "停止中")
	message.SetString(language.SimplifiedChinese, "Starting", "启动中")
	message.SetString(language.SimplifiedChinese, "Normal", "正常")
	message.SetString(language.SimplifiedChinese, "Stopped", "停止")
	message.SetString(language.SimplifiedChinese, "Creating", "创建中")
	message.SetString(language.SimplifiedChinese, "Deleting", "删除中")
	message.SetString(language.SimplifiedChinese, "Rebooting", "重启中")
	message.SetString(language.SimplifiedChinese, "DBInstanceClassChanging", "升降级中")
	message.SetString(language.SimplifiedChinese, "TRANSING", "迁移中")
	message.SetString(language.SimplifiedChinese, "EngineVersionUpgrading", "迁移版本中")
	message.SetString(language.SimplifiedChinese, "TransingToOthers	迁移数据到其他RDS", "中")
	message.SetString(language.SimplifiedChinese, "GuardDBInstanceCreating", "生产灾备实例中")
	message.SetString(language.SimplifiedChinese, "Restoring", "备份恢复中")
	message.SetString(language.SimplifiedChinese, "Importing", "数据导入中")
	message.SetString(language.SimplifiedChinese, "ImportingFromOthers	从其他RDS", "实例导入数据中")
	message.SetString(language.SimplifiedChinese, "DBInstanceNetTypeChanging", "内外网切换中")
	message.SetString(language.SimplifiedChinese, "GuardSwitching", "容灾切换中")
	message.SetString(language.SimplifiedChinese, "INS_CLONING", "实例克隆中")
	message.SetString(language.SimplifiedChinese, "Changing", "修改中")
	message.SetString(language.SimplifiedChinese, "Inactive", "被禁用")
	message.SetString(language.SimplifiedChinese, "Flushing", "清除中")
	message.SetString(language.SimplifiedChinese, "Released", "已释放")
	message.SetString(language.SimplifiedChinese, "Transforming", "转换中")
	message.SetString(language.SimplifiedChinese, "Unavailable", "服务停止")
	message.SetString(language.SimplifiedChinese, "Error", "创建失败")
	message.SetString(language.SimplifiedChinese, "Migrating", "迁移中")
	message.SetString(language.SimplifiedChinese, "BackupRecovering", "备份恢复中")
	message.SetString(language.SimplifiedChinese, "MinorVersionUpgrading", "小版本升级中")
	message.SetString(language.SimplifiedChinese, "NetworkModifying", "网络变更中")
	message.SetString(language.SimplifiedChinese, "SSLModifying", "SSL变更中")
	message.SetString(language.SimplifiedChinese, "MajorVersionUpgrading", "大版本升级中，可正常访问")
	message.SetString(language.SimplifiedChinese, "BeforeExpired (10 days)", "即将过期 (10天)")
	message.SetString(language.SimplifiedChinese, "Expired", "已过期")

	message.SetString(language.SimplifiedChinese, "running", "运行中")
	message.SetString(language.SimplifiedChinese, "stoped", "停止")
	message.SetString(language.SimplifiedChinese, "active", "正常")
	message.SetString(language.SimplifiedChinese, "activating", "生效中")
	message.SetString(language.SimplifiedChinese, "inactive", "冻结")
	message.SetString(language.SimplifiedChinese, "invalid", "失效")
	message.SetString(language.SimplifiedChinese, "locked", "锁定")
	// cloud resource ons
	message.SetString(language.SimplifiedChinese, "Standard edition instance", "标准版实例")
	message.SetString(language.SimplifiedChinese, "Platinum Edition instance", "铂金版实例")
	message.SetString(language.SimplifiedChinese, "Platinum Edition instance deploying", "铂金版实例部署中")
	message.SetString(language.SimplifiedChinese, "Standard edition instance expired", "标准版实例已欠费")
	message.SetString(language.SimplifiedChinese, "Upgrading", "升级中")
	// mysql spec
	message.SetString(language.SimplifiedChinese, "mysql.n1.micro.1", "基础版 1核1G")
	message.SetString(language.SimplifiedChinese, "mysql.n2.small.1", "基础版 1核2G")
	message.SetString(language.SimplifiedChinese, "mysql.n2.medium.1", "基础版 2核4G")
	message.SetString(language.SimplifiedChinese, "mysql.n2.large.1", "基础版 4核8G")
	message.SetString(language.SimplifiedChinese, "mysql.n2.xlarge.1", "基础版 8核16G")
	message.SetString(language.SimplifiedChinese, "mysql.n4.medium.1", "基础版 2核8G")
	message.SetString(language.SimplifiedChinese, "mysql.n4.large.1", "基础版 4核16G")
	message.SetString(language.SimplifiedChinese, "mysql.n4.xlarge.1", "基础版 8核32G")

	message.SetString(language.SimplifiedChinese, "rds.mysql.t1.small", "高可用版 通用型 1核1G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.s1.small", "高可用版 通用型 1核2G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.s2.large", "高可用版 通用型 2核4G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.s2.xlarge", "高可用版 通用型 2核8G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.s3.large", "高可用版 通用型 4核8G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.m1.medium", "高可用版 通用型 4核16G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.c1.large", "高可用版 通用型 8核16G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.c1.xlarge", "高可用版 通用型 8核32G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.c2.xlarge", "高可用版 通用型 16核64G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.c2.xlp2", "高可用版 通用型 16核96G")
	message.SetString(language.SimplifiedChinese, "rds.mysql.c2.2xlarge", "高可用版 通用型 16核128G")

	message.SetString(language.SimplifiedChinese, "mysql.x4.large.2", "高可用版 独享型	4核16G")
	message.SetString(language.SimplifiedChinese, "mysql.x4.xlarge.2", "高可用版 独享型 8核32G")
	message.SetString(language.SimplifiedChinese, "mysql.x4.2xlarge.2", "高可用版 独享型 16核64G")
	message.SetString(language.SimplifiedChinese, "mysql.x4.4xlarge.2", "高可用版 独享型 32核128G")
	message.SetString(language.SimplifiedChinese, "mysql.x8.medium.2", "高可用版 独享型 2核16G")
	message.SetString(language.SimplifiedChinese, "mysql.x8.large.2", "高可用版 独享型 4核32G")
	message.SetString(language.SimplifiedChinese, "mysql.x8.xlarge.2", "高可用版 独享型 8核64G")
	message.SetString(language.SimplifiedChinese, "mysql.x8.2xlarge.2", "高可用版 独享型 16核128G")
	message.SetString(language.SimplifiedChinese, "mysql.x8.4xlarge.2", "高可用版 独享型 32核256G")
	message.SetString(language.SimplifiedChinese, "mysql.x8.8xlarge.2", "高可用版 独享型 64核512G")

	message.SetString(language.SimplifiedChinese, "mysql.n2.small.2c", "高可用版 1核2G")
	message.SetString(language.SimplifiedChinese, "mysql.n2.medium.2c", "高可用版 2核4G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.medium.2c", "高可用版 2核4G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.large.2c", "高可用版 4核8G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.xlarge.2c", "高可用版 8核16G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.3large.2c", "高可用版 12核24G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.2xlarge.2c", "高可用版 16核32G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.3xlarge.2c", "高可用版 24核48G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.4xlarge.2c", "高可用版 32核64G")
	message.SetString(language.SimplifiedChinese, "mysql.x2.13large.2c", "高可用版 52核96G")

	// redis spec
	message.SetString(language.SimplifiedChinese, "redis.master.micro.default", "主从版 256MB")
	message.SetString(language.SimplifiedChinese, "redis.master.small.default", "主从版 1G")
	message.SetString(language.SimplifiedChinese, "redis.master.mid.default", "主从版 2G")
	message.SetString(language.SimplifiedChinese, "redis.master.stand.default", "主从版 4G")
	message.SetString(language.SimplifiedChinese, "redis.master.large.default", "主从版 8G")
	message.SetString(language.SimplifiedChinese, "redis.master.2xlarge.default", "主从版 16G")
	message.SetString(language.SimplifiedChinese, "redis.master.4xlarge.default", "主从版 32G")
	message.SetString(language.SimplifiedChinese, "redis.master.8xlarge.default", "主从版 64G")

	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.512m.2db.0rodb.4proxy.default", "集群版 1G (2节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.1g.2db.0rodb.4proxy.default", "集群版 2G (2节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.2g.2db.0rodb.4proxy.default", "集群版 4G (2节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.4g.2db.0rodb.4proxy.default", "集群版 8G (2节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.8g.2db.0rodb.4proxy.default", "集群版 16G (2节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.16g.2db.0rodb.4proxy.default", "集群版 32G (2节点)")

	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.512m.4db.0rodb.4proxy.default", "集群版 2G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.1g.4db.0rodb.4proxy.default", "集群版 4G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.2g.4db.0rodb.4proxy.default", "集群版 8G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.4g.4db.0rodb.4proxy.default", "集群版 16G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.6g.4db.0rodb.4proxy.default", "集群版 24G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.8g.4db.0rodb.4proxy.default", "集群版 32G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.10g.4db.0rodb.4proxy.default", "集群版 40G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.16g.4db.0rodb.4proxy.default", "集群版 64G (4节点)")
	message.SetString(language.SimplifiedChinese, "redis.logic.sharding.32g.4db.0rodb.8proxy.default", "集群版 128G (4节点)")

	message.SetString(language.SimplifiedChinese, "No cloud resource account is configured under the current org(%s)", "当前企业(%s)下没有配置云资源账号")

	message.SetString(language.SimplifiedChinese, "Shared Instance", "共享实例")

	message.SetString(language.SimplifiedChinese, "Basic Information", "基本信息")
	message.SetString(language.SimplifiedChinese, "Instance ID", "实例ID")
	message.SetString(language.SimplifiedChinese, "Region and Zone", "地域可用区")
	message.SetString(language.SimplifiedChinese, "Instance Role & Edition", "类型及系列")
	message.SetString(language.SimplifiedChinese, "Basic", "基础版")
	message.SetString(language.SimplifiedChinese, "HighAvailability", "高可用版")
	message.SetString(language.SimplifiedChinese, "AlwaysOn", "集群版")
	message.SetString(language.SimplifiedChinese, "Finance", "三节点企业版")
	message.SetString(language.SimplifiedChinese, "Storage Type", "存储类型")
	message.SetString(language.SimplifiedChinese, "local_ssd", "本地SSD盘")
	message.SetString(language.SimplifiedChinese, "ephemeral_ssd", "本地SSD盘")
	message.SetString(language.SimplifiedChinese, "cloud_ssd", "SSD云盘")
	message.SetString(language.SimplifiedChinese, "cloud_essd", "ESSD云盘")
	message.SetString(language.English, "local_ssd", "Local SSD")
	message.SetString(language.English, "ephemeral_ssd", "Local SSD")
	message.SetString(language.English, "cloud_ssd", "Cloud SSD")
	message.SetString(language.English, "cloud_essd", "Cloud ESSD")
	message.SetString(language.SimplifiedChinese, "Internal Endpoint", "内网地址")
	message.SetString(language.SimplifiedChinese, "Internal Port", "内网端口")
	message.SetString(language.SimplifiedChinese, "Public Endpoint", "外网地址")
	message.SetString(language.SimplifiedChinese, "Public Port", "外网端口")

	message.SetString(language.SimplifiedChinese, "Configuration Information", "配置信息")
	message.SetString(language.SimplifiedChinese, "Instance Family", "规格族")
	message.SetString(language.SimplifiedChinese, "s", "共享型")
	message.SetString(language.SimplifiedChinese, "x", "通用型")
	message.SetString(language.SimplifiedChinese, "d", "独享套餐")
	message.SetString(language.SimplifiedChinese, "h", "独占物理机")
	message.SetString(language.English, "s", "Shared-purpose")
	message.SetString(language.English, "x", "General-purpose")
	message.SetString(language.English, "d", "Dedicated Instance")
	message.SetString(language.English, "h", "Exclusive physical machine")

	message.SetString(language.SimplifiedChinese, "Database Engine", "数据库类型")
	message.SetString(language.SimplifiedChinese, "Memory", "数据库内存")
	message.SetString(language.SimplifiedChinese, "Maximum Connections", "最大连接数")
	message.SetString(language.SimplifiedChinese, "Maximum Iops", "最大IOPS")

	message.SetString(language.SimplifiedChinese, "Usage Statistics", "使用量统计")
	message.SetString(language.SimplifiedChinese, "Storage Capacity", "存储空间")
	message.SetString(language.SimplifiedChinese, "Used", "已使用")
	message.SetString(language.SimplifiedChinese, "Total", "共")

	message.SetString(language.SimplifiedChinese, "IndependentNaming", "独立命名空间")
	message.SetString(language.SimplifiedChinese, "Tcp Endpoint", "TCP 协议客户端接入点")
	message.SetString(language.SimplifiedChinese, "Http Endpoint", "HTTP 协议客户端接入点")
	message.SetString(language.SimplifiedChinese, "Private Host", "内网连接地址")
	message.SetString(language.SimplifiedChinese, "Public Host", "公网连接地址")
	message.SetString(language.SimplifiedChinese, "true", "是")
	message.SetString(language.SimplifiedChinese, "false", "否")

	// redis
	message.SetString(language.SimplifiedChinese, "Network Type", "网络类型")
	message.SetString(language.SimplifiedChinese, "CLASSIC", "经典网络")
	message.SetString(language.SimplifiedChinese, "VPC", "VPC网络")
	message.SetString(language.SimplifiedChinese, "VSwitch", "交换机信息")
	message.SetString(language.SimplifiedChinese, "MaxConnections", "最大连接数")
	message.SetString(language.SimplifiedChinese, "Connection Information", "连接信息")
	message.SetString(language.SimplifiedChinese, "Port", "端口")
	// ons
	message.SetString(language.SimplifiedChinese, "Normal message", "普通消息")
	message.SetString(language.SimplifiedChinese, "Partitionally ordered message", "分区顺序消息")
	message.SetString(language.SimplifiedChinese, "Globally ordered message", "全局顺序消息")
	message.SetString(language.SimplifiedChinese, "Transactional Message", "事务消息")
	message.SetString(language.SimplifiedChinese, "Scheduled/delayed message", "定时/延时消息")
	message.SetString(language.SimplifiedChinese, "owner", "持有者")
	message.SetString(language.SimplifiedChinese, "subscribable", "可以发布")
	message.SetString(language.SimplifiedChinese, "publishable", "可以订阅")
	message.SetString(language.SimplifiedChinese, "subscribable and publishable", "可以发布和订阅")

}
