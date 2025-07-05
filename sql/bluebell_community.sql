create table community
(
    id             int auto_increment
        primary key,
    community_id   int unsigned                        not null,
    community_name varchar(128)                        not null,
    introduction   varchar(256)                        not null,
    create_time    timestamp default CURRENT_TIMESTAMP not null,
    update_time    timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint idx_community_id
        unique (community_id),
    constraint idx_community_name
        unique (community_name)
)
    collate = utf8mb4_general_ci;

INSERT INTO bluebell.community (id, community_id, community_name, introduction, create_time, update_time) VALUES (1, 1, 'Go', 'Golang', '2025-06-07 08:30:00', '2025-06-07 08:30:00');
INSERT INTO bluebell.community (id, community_id, community_name, introduction, create_time, update_time) VALUES (2, 2, 'leetcode', '刷题刷题刷题', '2025-06-07 08:30:00', '2025-06-07 08:30:00');
INSERT INTO bluebell.community (id, community_id, community_name, introduction, create_time, update_time) VALUES (3, 3, '瓦罗兰特', '我不是瓦学弟', '2025-06-07 08:30:00', '2025-06-07 08:30:00');
INSERT INTO bluebell.community (id, community_id, community_name, introduction, create_time, update_time) VALUES (4, 4, '编程', '想成为高手', '2025-06-07 08:30:00', '2025-06-07 08:30:00');