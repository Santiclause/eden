delete from users;
delete from roles;
delete from permissions;
alter table users auto_increment = 1;
alter table roles auto_increment = 1;
alter table permissions auto_increment = 1;
