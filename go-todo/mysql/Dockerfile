FROM mysql:5.7.22

ENV MYSQL_ROOT_PASSWORD root
ADD ./init.sql /docker-entrypoint-initdb.d
ADD ./my.cnf /etc/mysql
