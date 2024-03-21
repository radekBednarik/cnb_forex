# Daily forex data app

Learning project using [SvelteKit](https://kit.svelte.dev), [Go](https://go.dev) and [Postgresql](https://www.postgresql.org)

## Preconditions

- Linux OS (Ubuntu based) - should work on Windows as well, but I did not try it
- Node LTS
- git
- postgresql
  - create a database `cnb_forex` with user, who have rights to create tables. You can use shell scripts in `/bin/db` folder. Modify them as you see fit.
- browser like Chromium or Firefox

## Installation

- clone via `git clone git@github.com:radekBednarik/cnb_forex.git`
- create a database and user using scripts in `/bin/db` folder
- switch to `/data-getter` folder and run `go mod tidy` and `go build main.go`
- then run `USER=<username> PASSWORD=<password> ./main` to download all data. Data are downloaded from public [Czech National Bank](https://www.cnb.cz) website starting from current date and going backwards to the date set in `/data-getter/config.toml` file. If you leave the setting, the download can take a while, since we are not using concurrency since we do not want to get banned on the server.
- switch to `/server` folder and run `go mod tidy` and `go build main.go`
- start server using `USER=<username> PASSWORD=<password> ./main`
- switch to `/frontend` and run `npm install`
- then run `npm run dev` to start dev server

App is ready to be used.
