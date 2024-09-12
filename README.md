docker run -d --name to_do_list --privileged=true -e MYSQL_ROOT_PASSWORD="admin" -e MYSQL_USER="thomas" -e MYSQL_PASSWORD="12345678" -e MYSQL_DATABASE="to_do_list" -p 3309:3306 mysql:latest
