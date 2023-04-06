DB
- docker image `docker pull postgres:15.2-alpine`
- docker run `docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.2-alpine`
- Install any macos postgres client like pgadmin
- migration
  - Install go migrate CLI
    - `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
    - `export PATH=$PATH:/Users/mahesh.singh/go/bin' >> ~/.zshrc` or `export PATH=$PATH:/Users/mahesh.singh/go/bin' >> ~/.bash_profile`
    - check `migrate -version`
  - to create a new migration
    - `migrate create -ext sql -dir db/migration -seq <migration title/file name>`
  - Install sqlc
    - `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`
    - how to use sqlc
      - Create CRUD queries in ./db/query directory with .sql via follwoing the syntax from https://docs.sqlc.dev/en/stable/howto/select.html
      - 