CREATE DATABASE IF NOT EXISTS eden;
ALTER DATABASE eden CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE IF NOT EXISTS users (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `username` varchar(60) NOT NULL,
    `password` CHAR(60) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL,
	`created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS roles (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(60) NOT NULL
);
CREATE TABLE IF NOT EXISTS permissions (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(60) NOT NULL
);
CREATE TABLE IF NOT EXISTS user_roles (
    `user_id` bigint NOT NULL,
    `role_id` bigint NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);
CREATE TABLE IF NOT EXISTS role_permissions (
    `role_id` bigint NOT NULL,
    `permission_id` bigint NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);
CREATE TABLE IF NOT EXISTS ircUsers (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `user_id` bigint NOT NULL,
    `nickname` varchar(60) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(id),
    UNIQUE KEY (`nickname`)
);
CREATE TABLE IF NOT EXISTS songs (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `filename` varchar(512),
    `title` varchar(255) NOT NULL,
    `artist_id` bigint NOT NULL,
    -- some other shit here
    FOREIGN KEY (artist_id) REFERENCES artists(id)
);
CREATE TABLE IF NOT EXISTS artists (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    UNIQUE KEY (name)
);
CREATE TABLE IF NOT EXISTS faves (
    `user_id` bigint NOT NULL,
    `song_id` bigint NOT NULL,
    PRIMARY KEY (`user_id`, `song_id`),
    FOREIGN KEY (`user_id`) REFERENCES users(id),
    FOREIGN KEY (`song_id`) REFERENCES songs(id)
);
CREATE TABLE IF NOT EXISTS playHistory (
    `id` bigint PRIMARY KEY AUTO_INCREMENT,
    `song_id` bigint NOT NULL,
	`dj` bigint NOT NULL,
    `played` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY (`played`),
    FOREIGN KEY (`song_id`) REFERENCES songs(id)
	FOREIGN KEY (`dj`) reference users(id)
);
