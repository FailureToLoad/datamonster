FROM postgres
ARG SA_PASSWORD
ENV POSTGRES_PASSWORD=${SA_PASSWORD:-Password1}
ENV POSTGRES_USER=admin

COPY init-db.sh /docker-entrypoint-initdb.d/
