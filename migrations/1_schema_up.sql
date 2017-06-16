CREATE DATABASE IF NOT EXISTS eden;
ALTER DATABASE eden CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE IF NOT EXISTS users (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `username` varchar(60) NOT NULL,
    `password` CHAR(60) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL
);
CREATE TABLE IF NOT EXISTS roles (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(60) NOT NULL
);
CREATE TABLE IF NOT EXISTS permissions (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(60) NOT NULL
);
CREATE TABLE IF NOT EXISTS rolePermissions (
    `role` bigint NOT NULL,
    `permission` bigint NOT NULL,
    PRIMARY KEY (role, permission),
    FOREIGN KEY (role) REFERENCES roles(id),
    FOREIGN KEY (permission) REFERENCES permissions(id)
);
CREATE TABLE IF NOT EXISTS songs (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `filename` varchar(512),
    `title` varchar(255) NOT NULL,
    `artist` bigint NOT NULL,
    -- some other shit here
    FOREIGN KEY (artist) REFERENCES artists(id)
);
CREATE TABLE IF NOT EXISTS artists (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    UNIQUE KEY (name)
);
CREATE TABLE IF NOT EXISTS faves (
    `user` bigint NOT NULL,
    `song` bigint NOT NULL,
    PRIMARY KEY (`user`, `song`),
    FOREIGN KEY (`user`) REFERENCES users(id),
    FOREIGN KEY (`song`) REFERENCES songs(id)
);
CREATE TABLE IF NOT EXISTS playHistory (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `song` bigint NOT NULL,
    `played` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY (`played`),
    FOREIGN KEY (`song`) REFERENCES songs(id)
);
