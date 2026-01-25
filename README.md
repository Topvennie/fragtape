# Fragtape

Fragtape is a self-hosted tool for **automatically generating CS2 highlight clips** and collecting fun group statistics.
It does **not** aim to be a performance analysis tool (for that I highly recommend [Leetify](https://leetify.com)).

> [!IMPORTANT]
> **Work in Progress**
>
> Fragtape is under active development and **not feature-complete**.
> Expect breaking changes and missing features.

Planned features include:

- Automatic detection of new CS2 matches
- Faceit integration
- Demo parsing to extract highlight moments
- Configurable highlight rules (kills, clutches, intervals, ...)
- Rendering clips using CS2 demos
- Group based settings and actions (e.g. Discord posting)
- Simple group stats (knife kills, team kills, ...)

## High Level Architecture

- **Server**
  - User-facing API
  - Stores configuration and metadata
  - Serves highlight clips
- **Worker**
  - Detects new matches / downloads demos
  - Parses demos and plans highlight segments
  - Creates recording jobs
  - Processes completed clips (e.g. Discord posting)
- **Recorder**
  - Consumes recording jobs
  - Runs CS2 + tooling to generate clips
  - Uploads results to object storage

## Recorder (Windows VM)

> [!IMPORTANT]
>
> This section is required for both **development** and **production** because the recorder container is part of the stack.

Fragtape’s recorder runs on Windows because the demo-recording tooling [HLAE](https://github.com/advancedfx/advancedfx) is Windows-only.

To make self-hosting easier, Fragtape uses **dockurr/windows**, which runs a Windows VM inside a Docker container.

There are two modes:

- **Development (dummy mode):** no CS2 / no HLAE required. Recorder produces dummy clips.
- **Real rendering:** requires GPU passthrough + CS2 + HLAE.

In both cases the majority of the setup is required.

### Configure VM resources

Set the resources for the Windows VM. You can do this in your `docker-compose.yml` or `.env`.
If you plan on actually recording the highlights be sure to check the minimum requirements for CS2.

- `RECORDER_DISK_SIZE` (example: `128G`)
- `RECORDER_RAM_SIZE` (example: `8G`)
- `RECORDER_CPU_CORES` (example: `8`)

See `.env.example`.

### First boot: install Windows

Start the recorder once before the rest of the stack:

```bash
docker compose up recorder
```

You can follow the download progress in the container logs.
Once the download has completed and the installation process starts you can switch over to `http://localhost:8006` for the progress.

### Connect to the VM

Once the download and installation process are complete, I highly recommend using RDP instead of the web UI for a much better experience.

For example to connect with [FreeRDP](https://www.freerdp.com/)

```bash
xfreerdp3 /v:127.0.0.1 /u:fragtape /p:admin
```

### Configure recorder connectivity (DB + MinIO)

The Windows VM must be able to react:

- Postgres
- MinIO

You can do this by changing the [config](./config/production.yml) values.
Get the ip of your docker host with `ip a` and change the following values:

- `recorder.db.host` -> `<ip>`
- `recorder.minio.endpoint` -> `<ip>:9000`

### Time and timezone (MinIO compatibility)

MinIO signing is sensitive to clock skew.

If you see errors like:

> The difference between the request time and the server's time is too large.
> your Windows clock/timezone is not in sync.

Fix by ensuring Windows time is correct and consistent with the host/containers.
Using UTC is recommended.

If you haven't changed the docker compose file then this is all correct.

### Rendering

This section is only relevant if set the config key `recorder.dummy_data` to false (this will try to run CS2 to render the highlights).

#### GPU passthrough

Passing a GPU through to dockurr VM is highly host / GPU specific.
I recommend searching for your specific setup and having look at the [Windows VM repository](https://github.com/dockur/windows).

To change the docker compose file I recommend to use a override file.
Docker compose will automatically load a file named `docker-compose.override.yml` if it exists.

<details>
  <summary>
    As an example, the override file I use for my ubuntu laptop.
  </summary>

```
services:
  recorder:
    environment:
      ARGUMENTS: >-
        -device vfio-pci,host=0000:01:00.0,multifunction=on
        -device vfio-pci,host=0000:01:00.1

    devices:
      - /dev/vfio:/dev/vfio

    cap_add:
      IPC_LOCK

    ulimits:
      memlock:
        soft: -1
        hard: -1

    security_opt:
      - seccomp=unconfined
```

</details>

After booting the VM, confirm the GPU appears in Windows Task Manager.

#### Install drivers, CS2 and HLAE

Inside the VM install:

- GPU drivers
- CS2
- [HLAE](https://github.com/advancedfx/advancedfx)

## Production

1. Copy `docker-compose.prod.yml` -> `docker-compose.yml`.
2. Copy `.env.example` -> `.env`.
3. Follow the recorder (Windows VM) section above.
4. Run `docker compose up -d`.
5. The server is reachable on port **8000**.

To update:

```bash
docker compose pull
docker compose down
docker compose up -d
```

## Development

### Quick Start

1. Install the tools listed in the [asdf file](./.tool-versions)
2. Install _make_.
3. Run `make setup` to install:

- Golang tools: _Air_, _Goose_, _Sqlc_, _Deadcode_
- Frontend dependencies

4. Install the git hook for code quality: `git config --local core.hooksPath .githooks/`
5. Copy `.env.example` -> `.env` and populate
6. Follow the recorder section in the readme
7. Run database migrations: `make migrate`
8. Start the project `make watch`.

Endpoints:

- **Backend:** <http://localhost:3001>
- **Frontend:** <http://localhost:3000>

By default you don't need to change any config values.
In any case don't change the log file of the recorder as the manager uses it to show the logs in the console.

