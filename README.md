docker run -d --name social_todo_list --privileged=true -e MYSQL_ROOT_PASSWORD="admin" -e MYSQL_USER="thomas" -e MYSQL_PASSWORD="12345678" -e MYSQL_DATABASE="social_todo_list" -p 3309:3306 mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
