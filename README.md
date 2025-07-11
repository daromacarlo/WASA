# WasaText

[![Università di Roma La Sapienza](https://img.shields.io/badge/university-La%20Sapienza-maroon)](https://www.uniroma1.it/)
[![Corso](https://img.shields.io/badge/course-Web%20and%20Software%20Architecture-orange)](https://corsidilaurea.uniroma1.it/it/view-course-details/2023/30027/20230113131042/c5bbbcfa-2298-4182-b66d-b625fc525307/7d8a1191-6d13-47a5-a2cd-eecb1876c238/3d9e27c4-5329-411f-85b1-46a0636e0934)

Questo progetto è pubblico e condiviso con l'intento di essere utile a chiunque voglia modificarlo e riutilizzarlo per il proprio lavoro o per fini didattici. Sei libero di copiarlo, distribuirlo e adattarlo come preferisci (tanto io l'esame l'ho fatto).

Questi sono i voti parziali del mio progetto:

![immagine](https://github.com/user-attachments/assets/e0a25378-b3a6-4668-9cd6-e3d380cdb6c8)

<sub><i>last update 10/07/2025</i></sub>

# Fantastic decaffeinated coffee
[![Go](https://img.shields.io/badge/language-Go-00ADD8)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/frontend-Vue.js-41B883)](https://vuejs.org/)
[![Yarn](https://img.shields.io/badge/package%20manager-Yarn-2C8EBB)](https://yarnpkg.com/)
[![OpenAPI](https://img.shields.io/badge/API-OpenAPI%203.0-yellow)](https://swagger.io/specification/)
[![Docker](https://img.shields.io/badge/container-Docker-2496ED)](https://www.docker.com/)
## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is an example of a daemon for checking the health of servers daemons; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the documentation (usually, for APIs, this means an OpenAPI file)
* `service/` has all packages for implementing project-specific functionalities
	* `service/api` contains an example of an API server
	* `service/globaltime` contains a wrapper package for `time.Time` (useful in unit testing)
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is an example of a web frontend in Vue.js; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

Other project files include:
* `open-node.sh` starts a new (temporary) container using `node:20` image for safe and secure web frontend development (you don't want to use `node` in your system, do you?).

## Go vendoring

This project uses [Go Vendoring](https://go.dev/ref/mod#vendoring). You must use `go mod vendor` after changing some dependency (`go get` or `go mod tidy`) and add all files under `vendor/` directory in your commit.

For more information about vendoring:

* https://go.dev/ref/mod#vendoring
* https://www.ardanlabs.com/blog/2020/04/modules-06-vendoring.html

## Node/YARN vendoring

This repository uses `yarn` and a vendoring technique that exploits the ["Offline mirror"](https://yarnpkg.com/features/caching). As for the Go vendoring, the dependencies are inside the repository.

You should commit the files inside the `.yarn` directory.

## How to set up a new project from this template

You need to:

* Change the Go module path to your module path in `go.mod`, `go.sum`, and in `*.go` files around the project
* Rewrite the API documentation `doc/api.yaml`
* If no web frontend is expected, remove `webui` and `cmd/webapi/register-webui.go`
* Update top/package comment inside `cmd/webapi/main.go` to reflect the actual project usage, goal, and general info
* Update the code in `run()` function (`cmd/webapi/main.go`) to connect to databases or external resources
* Write API code inside `service/api`, and create any further package inside `service/` (or subdirectories)

## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```shell
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-node.sh
# (here you're inside the container)
yarn run build-embed
exit
# (outside the container)
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

You can launch the backend only using:

```shell
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

```shell
./open-node.sh
# (here you're inside the container)
yarn run dev
```

## How to build for production / homework delivery

```shell
./open-node.sh
# (here you're inside the container)
yarn run build-prod
```

For "Web and Software Architecture" students: before committing and pushing your work for grading, please read the section below named "My build works when I use `yarn run dev`, however there is a Javascript crash in production/grading"

## Known issues

### My build works when I use `yarn run dev`, however there is a Javascript crash in production/grading

Some errors in the code are somehow not shown in `vite` development mode. To preview the code that will be used in production/grading settings, use the following commands:

```shell
./open-node.sh
# (here you're inside the container)
yarn run build-prod
yarn run preview
```

## License

See [LICENSE](LICENSE).
