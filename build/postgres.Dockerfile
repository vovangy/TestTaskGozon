FROM postgres:latest

RUN apt-get update && apt-get install -y postgresql-contrib

#FROM migrate/migrate:v4.15.2
#COPY ./db/migrations /migrations

#ENTRYPOINT ["sh", "-c"]
#CMD ["-path", "/migrations", "postgres://$USERNAME:$PASSWORD@$POSTGRES_HOST:5432/$DATABASE_NAME?sslmode=disable", "up"]