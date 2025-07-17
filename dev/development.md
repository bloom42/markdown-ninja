# Developement


## Docker

<!-- Setup Docker for multi-platform images builds:
```shell
docker buildx create --name docker_builder_multiplatform --bootstrap --use
``` -->

Create a `no-internet` network to block telemetry.

```bash
$ docker network create --internal --subnet 172.19.0.0/16 no-internet
```

resources:
  - https://docs.docker.com/build/building/multi-platform


## Database

Run the `timescale/timescaledb` container:

```bash
$ docker run --name markdown_ninja_db -d --network no-internet -e TIMESCALEDB_TELEMETRY=off -e POSTGRES_USER=markdown_ninja -e POSTGRES_PASSWORD=[PASSWORD] timescale/timescaledb:latest-pg16
```

resources:
  - https://docs.timescale.com/self-hosted/latest/install/installation-docker/


## Devcontainer

Connect your developemnt container to the `no-internet` network:

```bash
$ docker network connect no-internet [YOUR_DEV_CONTAINER]
```


## Stripe

```bash
docker run --rm -it -e STRIPE_API_KEY="XXX" -e STRIPE_DEVICE_NAME="dev" -e STRIPE_CLI_TELEMETRY_OPTOUT=1  stripe/stripe-cli listen --forward-to [docker_ip:docker_port]/api/webhooks/stripe --headers "Host:localhost"
```

The `Host` header is required to get requests correctly matched by the host router.
