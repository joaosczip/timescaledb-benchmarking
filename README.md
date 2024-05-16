# TimescaleDB Benchmarking

This app is a command-line tool to measure the performance of queries performed by a TimescaleDB instance.

## High-level overview and architecture

### Overview

The application uses the [cobra-cli](https://github.com/spf13/cobra) package to receive the input file through the command line.

Once the file is read, a buffered channel will be created based on the number of available CPU cores in the runner machine. And it'll launch a `worker pool` of goroutines to execute concurrently.

Each goroutine will do two things:
- Issue the query with the parameters from the input file against the database. One query for each line of the input file.
- Measure the query's execution time.

The data collected by each one of the goroutines will be gathered into the buffered channel.

After the process is finished, all the gathered information will be displayed to the user through the terminal.

### Architecture

![Untitled Diagram drawio (8)](https://github.com/joaosczip/timescaledb-benchmarking/assets/38441035/740f845f-405f-4d77-b087-8e2adf66c504)


## Requirements
- **go**: 1.20.2
- **Docker**: 20.10.17
- **Docker compose**: 2.21.0

You can check if your environment meets the requirements above by running:

```sh
$ go version
go version go1.20.2

$ docker -v
Docker version 20.10.17, build 100c701

# it must be docker compose, and not docker-compose (the old one)
$ docker compose version
Docker Compose version v2.21.0
```

## Running the tests

To run the automated tests, run the command specified below:

```sh
$ go test ./...
```

## Set up the environment
The `docker-compose.yml` file will build both the containers for the application and the database.

The application is a simple `golang` command line app that receives a path to a file containing a bunch of `query parameters` that will be send to the database.

The database is a `postgresql` instance with the `TimescaleDB` extension installed.

In order to have both of them running, you may use the `docker compose` tool:

```sh
$ docker compose up -d --build
```

This command will build and start both of the containers. 

**Remember to have the `docker compose` tool installed rather than the `docker-compose`.**

If everything went successfully, you should see the following in your terminal:

```sh
[+] Running 2/2
 ✔ Container db               Started                                                                                                                                       0.1s 
 ✔ Container timescale-app-1  Started 
```

You may also execute the following:

```sh
$ docker ps

# then you should receive an output like this one
CONTAINER ID   IMAGE                               COMMAND                  CREATED          STATUS                 PORTS                                                                                                           NAMES
feb324beed61   timescale-app                       "tail -f /dev/null"      49 seconds ago   Up 37 seconds                                                                                                                          timescale-app-1
acce8998cf21   timescale/timescaledb:latest-pg14   "docker-entrypoint.s…"   49 seconds ago   Up 37 seconds          0.0.0.0:5432->5432/tcp, :::5432->5432/tcp
```

## Running the benchmark
Once both of the containers are up and running, execute the following command to start it:

```sh
$ docker compose exec app ./main queryMetrics -p configs/query_params.csv
```

Then you should receive an output like the following:

```sh
+----------------+-------------------+------------+---------------+----------+----------------+----------------+
| AVG QUERY TIME | MEDIAN QUERY TIME | TOTAL TIME | TOTAL QUERIES | FAILURES | MIN QUERY TIME | MAX QUERY TIME |
+----------------+-------------------+------------+---------------+----------+----------------+----------------+
| 0.53ms         | 0.56ms            | 106.90ms   |           200 |        0 | 0.33ms         | 0.64ms         |
+----------------+-------------------+------------+---------------+----------+----------------+----------------+
```
