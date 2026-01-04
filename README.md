# Datamonster

Datamonster is a hobby project I've been working on, restarting, and allowing to live rent free in my head for a few years now. This repo is the current incarnation of it.

Kingdom Death: Monster is a board game with a _lot_ of book keeping. The game is played over a series of "lantern years" that make up a campaign's timeline. Each lantern year is split into three phases - develop, hunt, and showdown. My partner and I have been playing for years and on a good night we can clear a lantern year in a little under an hour.

Thing is, a campaign is comprised of anywhere between 20-30 lantern years.

I truly believe automating as much of the phase to phase play will streamline our ability to play this game.

## Running Locally

### Requirements and Setup

This project has the following requirements.

- Go 1.25.5 or higher
- Typescript 5.9.3
- node 24.10.1
- pnpm 10.24.0
- podman 5.7.1

The project relies on standing up an authelia container, a valkey container, a postgres container, and a small file server called glossary. The authelia container requires trusting self-signed local certs in order to function. If you're on linux, the process looks like this.

```bash
make certs
make valkey
make authelia
make postgres
make glossary
```

As part of setup for the frontend, make sure to also run the command to pull daisyui styles. I've complicated things a wee bit by using daisyui without tailwind, but I really enjoy the project more this way.

Below is something resembling the path you'll have to take to get dependencies installed.

```bash
cd api
go mod tidy
cd ../web
pnpm i
make daisyui
```

If you're using podman you'll also need to set these environment variables to ensure the Go tests can run without issues.

```bash
export DOCKER_HOST=unix:///run/user/1000/podman/podman.sock
export TESTCONTAINERS_RYUK_DISABLED=true
```

For any other platforms, the `make certs` step is going to need some manual attention. Feel free to raise an issue in the repo if you have problems but otherwise I'll be keeping this focused on linux (Fedora) as my dev environment.

I don't use a compose for the containers because valkey and authelia are pretty hands off once they're stood up. I only really need to fuss with the postgres container every so often. The podman commands _should_ be interchangeable with docker so if you use docker you can edit the makefile locally.

The `launch.json` has all the environment variables set for running the api. There may be discrepencies in how your container runtime assigns ip addresses. The web project does not rely on any environment variables and the vite dev server already has reverse proxy configured for hitting the api.

### Running the Project

I generally run the api from vscode and run the front-end in a terminal with `pnpm dev`.  
The test user within authelia is `testuser` with a password of `Password1`.
