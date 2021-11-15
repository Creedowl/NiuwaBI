create table users
(
    id          bigint auto_increment comment '主键',
    created_at  timestamp                                                               not null default current_timestamp comment '创建时间',
    updated_at  timestamp                                                               not null default current_timestamp on update current_timestamp comment '更新时间',
    name        varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci default '' not null comment '用户名',
    nickname    varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci default '' not null comment '昵称',
    password    varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci default '' not null comment 'hashed 密码',
    permissions bigint                                                       default 0  not null comment '用户权限',
    constraint user_pk
        primary key (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE utf8mb4_unicode_ci comment '用户表';