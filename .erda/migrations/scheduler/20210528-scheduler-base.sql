-- MIGRATION_BASE

CREATE TABLE `s_instance_info`
(
    `id`                    bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`            timestamp NULL DEFAULT NULL,
    `updated_at`            timestamp NULL DEFAULT NULL,
    `cluster`               varchar(64)   DEFAULT NULL,
    `namespace`             varchar(64)   DEFAULT NULL,
    `name`                  varchar(64)   DEFAULT NULL,
    `org_name`              varchar(64)   DEFAULT NULL,
    `org_id`                varchar(64)   DEFAULT NULL,
    `project_name`          varchar(64)   DEFAULT NULL,
    `project_id`            varchar(64)   DEFAULT NULL,
    `application_name`      varchar(255)  DEFAULT NULL,
    `application_id`        varchar(255)  DEFAULT NULL,
    `runtime_name`          varchar(255)  DEFAULT NULL,
    `runtime_id`            varchar(255)  DEFAULT NULL,
    `service_name`          varchar(255)  DEFAULT NULL,
    `workspace`             varchar(10)   DEFAULT NULL,
    `service_type`          varchar(64)   DEFAULT NULL,
    `addon_id`              varchar(255)  DEFAULT NULL,
    `meta`                  varchar(255)  DEFAULT NULL,
    `task_id`               varchar(150)  DEFAULT NULL,
    `phase`                 varchar(255)  DEFAULT NULL,
    `message`               varchar(1024) DEFAULT NULL,
    `container_id`          varchar(100)  DEFAULT NULL,
    `container_ip`          varchar(255)  DEFAULT NULL,
    `host_ip`               varchar(255)  DEFAULT NULL,
    `started_at`            timestamp NULL DEFAULT NULL,
    `finished_at`           timestamp NULL DEFAULT NULL,
    `exit_code`             int(11) DEFAULT NULL,
    `cpu_origin`            double        DEFAULT NULL,
    `mem_origin`            int(11) DEFAULT NULL,
    `cpu_request`           double        DEFAULT NULL,
    `mem_request`           int(11) DEFAULT NULL,
    `cpu_limit`             double        DEFAULT NULL,
    `mem_limit`             int(11) DEFAULT NULL,
    `image`                 varchar(255)  DEFAULT NULL,
    `edge_application_name` varchar(255)  DEFAULT NULL COMMENT '边缘应用名',
    `edge_site`             varchar(255)  DEFAULT NULL COMMENT '边缘站点',
    PRIMARY KEY (`id`),
    KEY                     `idx_s_instance_info_namespace` (`namespace`),
    KEY                     `idx_s_instance_info_name` (`name`),
    KEY                     `idx_s_instance_info_org_name` (`org_name`),
    KEY                     `idx_s_instance_info_org_id` (`org_id`),
    KEY                     `idx_s_instance_info_project_id` (`project_id`),
    KEY                     `idx_s_instance_info_container_id` (`container_id`),
    KEY                     `idx_s_instance_info_cluster` (`cluster`),
    KEY                     `idx_s_instance_info_project_name` (`project_name`),
    KEY                     `idx_s_instance_info_task_id` (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='实例信息';

CREATE TABLE `s_pod_info`
(
    `id`               bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `created_at`       timestamp NULL DEFAULT NULL,
    `updated_at`       timestamp NULL DEFAULT NULL,
    `cluster`          varchar(64)   DEFAULT NULL,
    `namespace`        varchar(64)   DEFAULT NULL,
    `name`             varchar(64)   DEFAULT NULL,
    `org_name`         varchar(64)   DEFAULT NULL,
    `org_id`           varchar(64)   DEFAULT NULL,
    `project_name`     varchar(64)   DEFAULT NULL,
    `project_id`       varchar(64)   DEFAULT NULL,
    `application_name` varchar(255)  DEFAULT NULL,
    `application_id`   varchar(255)  DEFAULT NULL,
    `runtime_name`     varchar(255)  DEFAULT NULL,
    `runtime_id`       varchar(255)  DEFAULT NULL,
    `service_name`     varchar(255)  DEFAULT NULL,
    `workspace`        varchar(10)   DEFAULT NULL,
    `service_type`     varchar(64)   DEFAULT NULL,
    `addon_id`         varchar(255)  DEFAULT NULL,
    `uid`              varchar(128)  DEFAULT NULL,
    `k8s_namespace`    varchar(128)  DEFAULT NULL,
    `pod_name`         varchar(128)  DEFAULT NULL,
    `phase`            varchar(255)  DEFAULT NULL,
    `message`          varchar(1024) DEFAULT NULL,
    `pod_ip`           varchar(255)  DEFAULT NULL,
    `host_ip`          varchar(255)  DEFAULT NULL,
    `started_at`       timestamp NULL DEFAULT NULL,
    `cpu_request`      double        DEFAULT NULL,
    `mem_request`      int(11) DEFAULT NULL,
    `cpu_limit`        double        DEFAULT NULL,
    `mem_limit`        int(11) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY                `idx_s_pod_info_namespace` (`namespace`),
    KEY                `idx_s_pod_info_name` (`name`),
    KEY                `idx_s_pod_info_org_name` (`org_name`),
    KEY                `idx_s_pod_info_org_id` (`org_id`),
    KEY                `idx_s_pod_info_project_id` (`project_id`),
    KEY                `idx_s_pod_info_cluster` (`cluster`),
    KEY                `idx_s_pod_info_project_name` (`project_name`),
    KEY                `idx_s_pod_info_k8s_namespace` (`k8s_namespace`),
    KEY                `idx_s_pod_info_pod_name` (`pod_name`),
    KEY                `idx_s_pod_info_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='pod信息';

