insert users(username, password) values ("Santiclause", "$2y$10$ANDOMfrMgTP6y4sLGJKKseVXcYgFT2ipTJgECKGo4/YEx4fibPv5u"), ("Kittensmakemesmile", "$2y$10$ANDOMfrMgTP6y4sLGJKKseVXcYgFT2ipTJgECKGo4/YEx4fibPv5u"), ("kinky", "$2y$10$ANDOMfrMgTP6y4sLGJKKseVXcYgFT2ipTJgECKGo4/YEx4fibPv5u"), ("MrBrawl", "$2y$10$ANDOMfrMgTP6y4sLGJKKseVXcYgFT2ipTJgECKGo4/YEx4fibPv5u");
insert roles(name) values ("admin"), ("dj");
insert user_roles values (1,1), (2,1), (3,1), (4,2);
insert permissions(name) values("super"), ("topic"), ("dj");
insert role_permissions values (1,1), (1,2), (1,3), (2,2), (2,3);
