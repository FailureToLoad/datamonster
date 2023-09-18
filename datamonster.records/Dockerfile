FROM postgres
ARG SA_PASSWORD
ENV POSTGRES_PASSWORD=${SA_PASSWORD:-Password1}
ENV POSTGRES_USER=admin
ENV POSTGRES_DB=terms

COPY /seed /seed
COPY init-db-terms.sh /docker-entrypoint-initdb.d/
