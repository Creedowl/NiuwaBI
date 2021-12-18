create table workspaces
(
    id         bigint auto_increment comment '主键',
    created_at timestamp                                             not null default current_timestamp comment '创建时间',
    updated_at timestamp                                             not null default current_timestamp on update current_timestamp comment '更新时间',
    name       varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci   default '' not null comment '工作区名称',
    config     text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci not null comment '工作区配置',
    owner      bigint                                                         default 0 not null comment '创建人',
    constraint workspace_pk
        primary key (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE utf8mb4_unicode_ci comment '工作区表';
