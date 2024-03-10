FROM mysql:latest

ENV MYSQL_USER=julian
ENV MYSQL_PASSWORD=12345
ENV MYSQL_ROOT_PASSWORD=12345

WORKDIR /app

# Se copia el archivo de inicializacion de la database
COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 3306