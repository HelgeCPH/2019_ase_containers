We have a small application, a telephone book application. It consists of a [webserver](webserver/telbook_server.go) written in Go, which finds all `Person` records stored in a MongoDB [database](dbserver/db_setup.js) and returns an HTML page with those records.


# An example with Docker

## Building the DB Server
The following is written from my perspective, i.e. user `helgecph`, which you want to replace with yours.

```bash
$ cd dbserver/
$ docker build -t helgecph/dbserver .
$ cd ..
```

## Building the Webserver

```bash
$ cd webserver/
$ docker build -t helgecph/webserver .
$ cd ..
```

Now, check that both images are locally available.

```bash
$ docker images
REPOSITORY           TAG                 IMAGE ID            CREATED              SIZE
helgecph/webserver   latest              234d72f53176        7 seconds ago        826MB
helgecph/dbserver    latest              de160d51e0f0        About a minute ago   363MB
```

## Starting the Application Manually



```bash
$ mkdir $(pwd)/datadb
$ docker run -d -p 27017:27017 --name dbserver helgecph/dbserver
$ docker run -it -d --rm --name webserver --link dbserver -p 8080:8080 helgecph/webserver
```

Eventhough deprecated, on can `--link` the containers via the bridge network together.

```bash
$ docker ps -a
CONTAINER ID        IMAGE                COMMAND                  CREATED             STATUS              PORTS                      NAMES
0282fc8b2c41        helgecph/webserver   "/bin/sh -c ./telb..."   11 seconds ago      Up 10 seconds       0.0.0.0:8080->8080/tcp     webserver
06b85924f444        helgecph/dbserver    "docker-entrypoint..."   6 minutes ago       Up 6 minutes        0.0.0.0:27017->27017/tcp   dbserver
```

Calling the client on linked containers works as presented in the lecture:

```bash
$ docker run --rm --link webserver appropriate/curl:latest curl -s http://webserver:8080
<!DOCTYPE HTML>
<html>
    <head>
        <title>The Møllers</title>
    </head>
    <body>
        <h1>Telephone Book</h1>
        <hr>
        <table style="width:50%">
          <tr>
            <th>Index</th>
            <th>Name</th>
            <th>Phone</th>
            <th>Address</th>
            <th>City</th>
          </tr>

          <tr>
            <td>0</td>
            <td>Møller</td>
            <td>&#43;45 20 86 46 44</td>
            <td>Herningvej 8</td>
            <td>4800 Nykøbing F</td>
          </tr>

          <tr>
            <td>1</td>
            <td>A Egelund-Møller</td>
            <td>&#43;45 54 94 41 81</td>
            <td>Rønnebærparken 1 0011</td>
            <td>4983 Dannemare</td>
          </tr>

          <tr>
            <td>2</td>
            <td>A K Møller</td>
            <td>&#43;45 75 50 75 14</td>
            <td>Bregnerødvej 75, st. 0002</td>
            <td>3460 Birkerød</td>
          </tr>

          <tr>
            <td>3</td>
            <td>A Møller</td>
            <td>&#43;45 97 95 20 01</td>
            <td>Dalstræde 11 Heltborg</td>
            <td>7760 Hurup Thy</td>
          </tr>

        </table>
        <p></p>
        Data taken from <a href="https://www.krak.dk/person/resultat/møller">Krak.dk</a>
    </body>
</html>
```

```bash
$ docker stop webserver
$ docker stop dbserver
$ docker rm dbserver
```


Properly done, from now on on links containers via a shared network.

```bash
$ docker network create example-network
$ docker network ls
NETWORK ID          NAME                            DRIVER              SCOPE
d5a8f5d3b2c2        bridge                          bridge              local
9c9d24069da7        example-network                 bridge              local
bd11ae20c3ac        host                            host                local
51892d4cc44a        none                            null                local
$ docker run -d -p 27017:27017 --name dbserver --network=example-network helgecph/dbserver
$ docker run -it -d --rm --name webserver --network=example-network -p 8080:8080 helgecph/webserver
```


### Testing the Application

```bash
$ docker run --rm --network=example-network appropriate/curl:latest curl http://webserver:8080
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0<!DOCTYPE HTML>
<html>
    <head>
        <title>The Møllers</title>
    </head>
    <body>
        <h1>Telephone Book</h1>
        <hr>
        <table style="width:50%">
          <tr>
            <th>Index</th>
            <th>Name</th>
            <th>Phone</th>
            <th>Address</th>
            <th>City</th>
          </tr>

          <tr>
            <td>0</td>
            <td>Møller</td>
            <td>&#43;45 20 86 46 44</td>
            <td>Herningvej 8</td>
            <td>4800 Nykøbing F</td>
          </tr>

          <tr>
            <td>1</td>
            <td>A Egelund-Møller</td>
            <td>&#43;45 54 94 41 81</td>
            <td>Rønnebærparken 1 0011</td>
            <td>4983 Dannemare</td>
          </tr>

          <tr>
            <td>2</td>
            <td>A K Møller</td>
            <td>&#43;45 75 50 75 14</td>
            <td>Bregnerødvej 75, st. 0002</td>
            <td>3460 Birkerød</td>
          </tr>

          <tr>
            <td>3</td>
            <td>A Møller</td>
            <td>&#43;45 97 95 20 01</td>
            <td>Dalstræde 11 Heltborg</td>
            <td>7760 Hurup Thy</td>
          </tr>

        </table>
        <p></p>
        Data taken from <a href="https://www.krak.dk/person/resultat/møller">Krak.dk</a>
    </body>
</html>
100  1366  100  1366    0     0   443k      0 --:--:-- --:--:-- --:--:--  666k
```


## Stopping the Application Manually


```bash
$ docker stop dbserver
$ docker stop webserver
```

```bash
$ docker rm webserver
$ docker rm dbserver
```

## Starting the Application with Docker Compose


```yml
version: '3'
services:
  dbserver:
    image: helgecph/dbserver
    ports:
      - "27017:27017"
    networks:
      - outside

  webserver:
    image: helgecph/webserver
    ports:
      - "8080:8080"
    networks:
        - outside

  clidownload:
    image: appropriate/curl
    networks:
      - outside
    entrypoint: sh -c  "sleep 5 && curl http://webserver:8080"

networks:
  outside:
    external:
      name: example-network
```


```bash
$ docker-compose up
Creating example_webserver_1   ... done
Creating example_clidownload_1 ... done
Creating example_dbserver_1    ... done
Attaching to example_clidownload_1, example_dbserver_1, example_webserver_1
dbserver_1     | about to fork child process, waiting until server is ready for connections.
dbserver_1     | forked process: 26
dbserver_1     | 2019-10-27T18:09:54.439+0000 I  CONTROL  [main] ***** SERVER RESTARTED *****
dbserver_1     | 2019-10-27T18:09:54.442+0000 I  CONTROL  [main] Automatically disabling TLS 1.0, to force-enable TLS 1.0 specify --sslDisabledProtocols 'none'
dbserver_1     | 2019-10-27T18:09:54.446+0000 I  CONTROL  [initandlisten] MongoDB starting : pid=26 port=27017 dbpath=/data/db 64-bit host=7bc89b50923e
dbserver_1     | 2019-10-27T18:09:54.446+0000 I  CONTROL  [initandlisten] db version v4.2.1
dbserver_1     | 2019-10-27T18:09:54.446+0000 I  CONTROL  [initandlisten] git version: edf6d45851c0b9ee15548f0f847df141764a317e
dbserver_1     | 2019-10-27T18:09:54.446+0000 I  CONTROL  [initandlisten] OpenSSL version: OpenSSL 1.1.1  11 Sep 2018
dbserver_1     | 2019-10-27T18:09:54.447+0000 I  CONTROL  [initandlisten] allocator: tcmalloc
dbserver_1     | 2019-10-27T18:09:54.447+0000 I  CONTROL  [initandlisten] modules: none
dbserver_1     | 2019-10-27T18:09:54.447+0000 I  CONTROL  [initandlisten] build environment:
dbserver_1     | 2019-10-27T18:09:54.447+0000 I  CONTROL  [initandlisten]     distmod: ubuntu1804
dbserver_1     | 2019-10-27T18:09:54.447+0000 I  CONTROL  [initandlisten]     distarch: x86_64
dbserver_1     | 2019-10-27T18:09:54.447+0000 I  CONTROL  [initandlisten]     target_arch: x86_64
dbserver_1     | 2019-10-27T18:09:54.447+0000 I  CONTROL  [initandlisten] options: { net: { bindIp: "127.0.0.1", port: 27017, tls: { mode: "disabled" } }, processManagement: { fork: true, pidFilePath: "/tmp/docker-entrypoint-temp-mongod.pid" }, systemLog: { destination: "file", logAppend: true, path: "/proc/1/fd/1" } }
dbserver_1     | 2019-10-27T18:09:54.449+0000 I  STORAGE  [initandlisten]
dbserver_1     | 2019-10-27T18:09:54.449+0000 I  STORAGE  [initandlisten] ** WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine
dbserver_1     | 2019-10-27T18:09:54.449+0000 I  STORAGE  [initandlisten] **          See http://dochub.mongodb.org/core/prodnotes-filesystem
dbserver_1     | 2019-10-27T18:09:54.449+0000 I  STORAGE  [initandlisten] wiredtiger_open config: create,cache_size=487M,cache_overflow=(file_max=0M),session_max=33000,eviction=(threads_min=4,threads_max=4),config_base=false,statistics=(fast),log=(enabled=true,archive=true,path=journal,compressor=snappy),file_manager=(close_idle_time=100000,close_scan_interval=10,close_handle_minimum=250),statistics_log=(wait=0),verbose=[recovery_progress,checkpoint_progress],
dbserver_1     | 2019-10-27T18:09:55.062+0000 I  STORAGE  [initandlisten] WiredTiger message [1572199795:62083][26:0x7f708bb32b00], txn-recover: Set global recovery timestamp: (0,0)
dbserver_1     | 2019-10-27T18:09:55.071+0000 I  RECOVERY [initandlisten] WiredTiger recoveryTimestamp. Ts: Timestamp(0, 0)
dbserver_1     | 2019-10-27T18:09:55.080+0000 I  STORAGE  [initandlisten] Timestamp monitor starting
dbserver_1     | 2019-10-27T18:09:55.084+0000 I  CONTROL  [initandlisten]
dbserver_1     | 2019-10-27T18:09:55.084+0000 I  CONTROL  [initandlisten] ** WARNING: Access control is not enabled for the database.
dbserver_1     | 2019-10-27T18:09:55.084+0000 I  CONTROL  [initandlisten] **          Read and write access to data and configuration is unrestricted.
dbserver_1     | 2019-10-27T18:09:55.085+0000 I  CONTROL  [initandlisten]
dbserver_1     | 2019-10-27T18:09:55.087+0000 I  STORAGE  [initandlisten] createCollection: admin.system.version with provided UUID: fdd9ac7e-9486-421a-853f-afd420a068b8 and options: { uuid: UUID("fdd9ac7e-9486-421a-853f-afd420a068b8") }
dbserver_1     | 2019-10-27T18:09:55.097+0000 I  INDEX    [initandlisten] index build: done building index _id_ on ns admin.system.version
dbserver_1     | 2019-10-27T18:09:55.098+0000 I  SHARDING [initandlisten] Marking collection admin.system.version as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:55.098+0000 I  COMMAND  [initandlisten] setting featureCompatibilityVersion to 4.2
dbserver_1     | 2019-10-27T18:09:55.099+0000 I  SHARDING [initandlisten] Marking collection local.system.replset as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:55.100+0000 I  STORAGE  [initandlisten] Flow Control is enabled on this deployment.
dbserver_1     | 2019-10-27T18:09:55.100+0000 I  SHARDING [initandlisten] Marking collection admin.system.roles as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:55.101+0000 I  STORAGE  [initandlisten] createCollection: local.startup_log with generated UUID: 9e801b15-8e1a-4e0f-9857-5497f74f4371 and options: { capped: true, size: 10485760 }
dbserver_1     | 2019-10-27T18:09:55.110+0000 I  INDEX    [initandlisten] index build: done building index _id_ on ns local.startup_log
dbserver_1     | 2019-10-27T18:09:55.110+0000 I  SHARDING [initandlisten] Marking collection local.startup_log as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:55.112+0000 I  FTDC     [initandlisten] Initializing full-time diagnostic data capture with directory '/data/db/diagnostic.data'
dbserver_1     | 2019-10-27T18:09:55.114+0000 I  SHARDING [LogicalSessionCacheRefresh] Marking collection config.system.sessions as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:55.116+0000 I  STORAGE  [LogicalSessionCacheRefresh] createCollection: config.system.sessions with provided UUID: 22b361d2-7664-4c36-91f8-070e87b2f90c and options: { uuid: UUID("22b361d2-7664-4c36-91f8-070e87b2f90c") }
dbserver_1     | 2019-10-27T18:09:55.117+0000 I  NETWORK  [initandlisten] Listening on /tmp/mongodb-27017.sock
dbserver_1     | 2019-10-27T18:09:55.117+0000 I  NETWORK  [initandlisten] Listening on 127.0.0.1
dbserver_1     | 2019-10-27T18:09:55.118+0000 I  NETWORK  [initandlisten] waiting for connections on port 27017
dbserver_1     | child process started successfully, parent exiting
dbserver_1     | 2019-10-27T18:09:55.127+0000 I  INDEX    [LogicalSessionCacheRefresh] index build: done building index _id_ on ns config.system.sessions
dbserver_1     | 2019-10-27T18:09:55.136+0000 I  INDEX    [LogicalSessionCacheRefresh] index build: starting on config.system.sessions properties: { v: 2, key: { lastUse: 1 }, name: "lsidTTLIndex", ns: "config.system.sessions", expireAfterSeconds: 1800 } using method: Hybrid
dbserver_1     | 2019-10-27T18:09:55.136+0000 I  INDEX    [LogicalSessionCacheRefresh] build may temporarily use up to 500 megabytes of RAM
dbserver_1     | 2019-10-27T18:09:55.136+0000 I  INDEX    [LogicalSessionCacheRefresh] index build: collection scan done. scanned 0 total records in 0 seconds
dbserver_1     | 2019-10-27T18:09:55.137+0000 I  INDEX    [LogicalSessionCacheRefresh] index build: inserted 0 keys from external sorter into index in 0 seconds
dbserver_1     | 2019-10-27T18:09:55.140+0000 I  INDEX    [LogicalSessionCacheRefresh] index build: done building index lsidTTLIndex on ns config.system.sessions
dbserver_1     | 2019-10-27T18:09:55.142+0000 I  SHARDING [LogicalSessionCacheReap] Marking collection config.transactions as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:55.196+0000 I  NETWORK  [listener] connection accepted from 127.0.0.1:44292 #1 (1 connection now open)
dbserver_1     | 2019-10-27T18:09:55.197+0000 I  NETWORK  [conn1] received client metadata from 127.0.0.1:44292 conn1: { application: { name: "MongoDB Shell" }, driver: { name: "MongoDB Internal Client", version: "4.2.1" }, os: { type: "Linux", name: "Ubuntu", architecture: "x86_64", version: "18.04" } }
dbserver_1     |
dbserver_1     | 2019-10-27T18:09:55.205+0000 I  NETWORK  [conn1] end connection 127.0.0.1:44292 (0 connections now open)
dbserver_1     | /usr/local/bin/docker-entrypoint.sh: running /docker-entrypoint-initdb.d/db_setup.js
dbserver_1     | 2019-10-27T18:09:55.267+0000 I  NETWORK  [listener] connection accepted from 127.0.0.1:44294 #2 (1 connection now open)
dbserver_1     | 2019-10-27T18:09:55.268+0000 I  NETWORK  [conn2] received client metadata from 127.0.0.1:44294 conn2: { application: { name: "MongoDB Shell" }, driver: { name: "MongoDB Internal Client", version: "4.2.1" }, os: { type: "Linux", name: "Ubuntu", architecture: "x86_64", version: "18.04" } }
dbserver_1     | 2019-10-27T18:09:55.276+0000 I  SHARDING [conn2] Marking collection test.people as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:55.278+0000 I  STORAGE  [conn2] createCollection: test.people with generated UUID: 325d8d08-60c0-4f5e-ac64-db50020ac170 and options: {}
dbserver_1     | 2019-10-27T18:09:55.290+0000 I  INDEX    [conn2] index build: done building index _id_ on ns test.people
dbserver_1     | 2019-10-27T18:09:55.296+0000 I  NETWORK  [conn2] end connection 127.0.0.1:44294 (0 connections now open)
dbserver_1     |
dbserver_1     |
dbserver_1     | 2019-10-27T18:09:55.321+0000 I  CONTROL  [main] ***** SERVER RESTARTED *****
dbserver_1     | 2019-10-27T18:09:55.324+0000 I  CONTROL  [main] Automatically disabling TLS 1.0, to force-enable TLS 1.0 specify --sslDisabledProtocols 'none'
dbserver_1     | killing process with pid: 26
dbserver_1     | 2019-10-27T18:09:55.327+0000 I  CONTROL  [signalProcessingThread] got signal 15 (Terminated), will terminate after current cmd ends
dbserver_1     | 2019-10-27T18:09:55.328+0000 I  NETWORK  [signalProcessingThread] shutdown: going to close listening sockets...
dbserver_1     | 2019-10-27T18:09:55.328+0000 I  NETWORK  [signalProcessingThread] removing socket file: /tmp/mongodb-27017.sock
dbserver_1     | 2019-10-27T18:09:55.329+0000 I  -        [signalProcessingThread] Stopping further Flow Control ticket acquisitions.
dbserver_1     | 2019-10-27T18:09:55.330+0000 I  CONTROL  [signalProcessingThread] Shutting down free monitoring
dbserver_1     | 2019-10-27T18:09:55.330+0000 I  FTDC     [signalProcessingThread] Shutting down full-time diagnostic data capture
dbserver_1     | 2019-10-27T18:09:55.331+0000 I  STORAGE  [signalProcessingThread] Deregistering all the collections
dbserver_1     | 2019-10-27T18:09:55.331+0000 I  STORAGE  [signalProcessingThread] Timestamp monitor shutting down
dbserver_1     | 2019-10-27T18:09:55.332+0000 I  STORAGE  [signalProcessingThread] WiredTigerKVEngine shutting down
dbserver_1     | 2019-10-27T18:09:55.333+0000 I  STORAGE  [signalProcessingThread] Shutting down session sweeper thread
dbserver_1     | 2019-10-27T18:09:55.334+0000 I  STORAGE  [signalProcessingThread] Finished shutting down session sweeper thread
dbserver_1     | 2019-10-27T18:09:55.334+0000 I  STORAGE  [signalProcessingThread] Shutting down journal flusher thread
dbserver_1     | 2019-10-27T18:09:55.374+0000 I  STORAGE  [signalProcessingThread] Finished shutting down journal flusher thread
dbserver_1     | 2019-10-27T18:09:55.374+0000 I  STORAGE  [signalProcessingThread] Shutting down checkpoint thread
dbserver_1     | 2019-10-27T18:09:55.375+0000 I  STORAGE  [signalProcessingThread] Finished shutting down checkpoint thread
dbserver_1     | 2019-10-27T18:09:55.417+0000 I  STORAGE  [signalProcessingThread] shutdown: removing fs lock...
dbserver_1     | 2019-10-27T18:09:55.418+0000 I  CONTROL  [signalProcessingThread] now exiting
dbserver_1     | 2019-10-27T18:09:55.419+0000 I  CONTROL  [signalProcessingThread] shutting down with code:0
dbserver_1     |
dbserver_1     | MongoDB init process complete; ready for start up.
dbserver_1     |
dbserver_1     | 2019-10-27T18:09:56.410+0000 I  CONTROL  [main] Automatically disabling TLS 1.0, to force-enable TLS 1.0 specify --sslDisabledProtocols 'none'
dbserver_1     | 2019-10-27T18:09:56.415+0000 I  CONTROL  [initandlisten] MongoDB starting : pid=1 port=27017 dbpath=/data/db 64-bit host=7bc89b50923e
dbserver_1     | 2019-10-27T18:09:56.415+0000 I  CONTROL  [initandlisten] db version v4.2.1
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten] git version: edf6d45851c0b9ee15548f0f847df141764a317e
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten] OpenSSL version: OpenSSL 1.1.1  11 Sep 2018
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten] allocator: tcmalloc
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten] modules: none
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten] build environment:
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten]     distmod: ubuntu1804
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten]     distarch: x86_64
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten]     target_arch: x86_64
dbserver_1     | 2019-10-27T18:09:56.416+0000 I  CONTROL  [initandlisten] options: { net: { bindIp: "*" } }
dbserver_1     | 2019-10-27T18:09:56.418+0000 I  STORAGE  [initandlisten] Detected data files in /data/db created by the 'wiredTiger' storage engine, so setting the active storage engine to 'wiredTiger'.
dbserver_1     | 2019-10-27T18:09:56.419+0000 I  STORAGE  [initandlisten]
dbserver_1     | 2019-10-27T18:09:56.419+0000 I  STORAGE  [initandlisten] ** WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine
dbserver_1     | 2019-10-27T18:09:56.419+0000 I  STORAGE  [initandlisten] **          See http://dochub.mongodb.org/core/prodnotes-filesystem
dbserver_1     | 2019-10-27T18:09:56.420+0000 I  STORAGE  [initandlisten] wiredtiger_open config: create,cache_size=487M,cache_overflow=(file_max=0M),session_max=33000,eviction=(threads_min=4,threads_max=4),config_base=false,statistics=(fast),log=(enabled=true,archive=true,path=journal,compressor=snappy),file_manager=(close_idle_time=100000,close_scan_interval=10,close_handle_minimum=250),statistics_log=(wait=0),verbose=[recovery_progress,checkpoint_progress],
dbserver_1     | 2019-10-27T18:09:57.033+0000 I  STORAGE  [initandlisten] WiredTiger message [1572199797:33784][1:0x7f056ec0fb00], txn-recover: Recovering log 1 through 2
dbserver_1     | 2019-10-27T18:09:57.168+0000 I  STORAGE  [initandlisten] WiredTiger message [1572199797:168591][1:0x7f056ec0fb00], txn-recover: Recovering log 2 through 2
dbserver_1     | 2019-10-27T18:09:57.482+0000 I  STORAGE  [initandlisten] WiredTiger message [1572199797:482156][1:0x7f056ec0fb00], txn-recover: Main recovery loop: starting at 1/27520 to 2/256
dbserver_1     | 2019-10-27T18:09:57.628+0000 I  STORAGE  [initandlisten] WiredTiger message [1572199797:628066][1:0x7f056ec0fb00], txn-recover: Recovering log 1 through 2
dbserver_1     | 2019-10-27T18:09:57.743+0000 I  STORAGE  [initandlisten] WiredTiger message [1572199797:743042][1:0x7f056ec0fb00], txn-recover: Recovering log 2 through 2
dbserver_1     | 2019-10-27T18:09:57.806+0000 I  STORAGE  [initandlisten] WiredTiger message [1572199797:806425][1:0x7f056ec0fb00], txn-recover: Set global recovery timestamp: (0,0)
dbserver_1     | 2019-10-27T18:09:57.819+0000 I  RECOVERY [initandlisten] WiredTiger recoveryTimestamp. Ts: Timestamp(0, 0)
dbserver_1     | 2019-10-27T18:09:57.825+0000 I  STORAGE  [initandlisten] Timestamp monitor starting
dbserver_1     | 2019-10-27T18:09:57.827+0000 I  CONTROL  [initandlisten]
dbserver_1     | 2019-10-27T18:09:57.828+0000 I  CONTROL  [initandlisten] ** WARNING: Access control is not enabled for the database.
dbserver_1     | 2019-10-27T18:09:57.828+0000 I  CONTROL  [initandlisten] **          Read and write access to data and configuration is unrestricted.
dbserver_1     | 2019-10-27T18:09:57.828+0000 I  CONTROL  [initandlisten]
dbserver_1     | 2019-10-27T18:09:57.833+0000 I  SHARDING [initandlisten] Marking collection local.system.replset as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:57.834+0000 I  STORAGE  [initandlisten] Flow Control is enabled on this deployment.
dbserver_1     | 2019-10-27T18:09:57.835+0000 I  SHARDING [initandlisten] Marking collection admin.system.roles as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:57.835+0000 I  SHARDING [initandlisten] Marking collection admin.system.version as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:57.837+0000 I  SHARDING [initandlisten] Marking collection local.startup_log as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:57.838+0000 I  FTDC     [initandlisten] Initializing full-time diagnostic data capture with directory '/data/db/diagnostic.data'
dbserver_1     | 2019-10-27T18:09:57.839+0000 I  SHARDING [LogicalSessionCacheReap] Marking collection config.system.sessions as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:57.839+0000 I  NETWORK  [initandlisten] Listening on /tmp/mongodb-27017.sock
dbserver_1     | 2019-10-27T18:09:57.840+0000 I  NETWORK  [initandlisten] Listening on 0.0.0.0
dbserver_1     | 2019-10-27T18:09:57.840+0000 I  SHARDING [LogicalSessionCacheReap] Marking collection config.transactions as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:57.840+0000 I  NETWORK  [initandlisten] waiting for connections on port 27017
dbserver_1     | 2019-10-27T18:09:58.001+0000 I  SHARDING [ftdc] Marking collection local.oplog.rs as collection version: <unsharded>
clidownload_1  |   % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
clidownload_1  |                                  Dload  Upload   Total   Spent    Left  Speed
dbserver_1     | 2019-10-27T18:09:59.301+0000 I  NETWORK  [listener] connection accepted from 172.19.0.3:51326 #1 (1 connection now open)
dbserver_1     | 2019-10-27T18:09:59.306+0000 I  SHARDING [conn1] Marking collection test.people as collection version: <unsharded>
dbserver_1     | 2019-10-27T18:09:59.307+0000 I  NETWORK  [conn1] end connection 172.19.0.3:51326 (0 connections now open)
100  1366  100  1366    0  <!DOCTYPE HTML>  0 --:--:-- --:--:-- --:--:--     0
clidownload_1  | <html>
clidownload_1  |     <head>
clidownload_1  |         <title>The Møllers</title>
clidownload_1  |     </head>
clidownload_1  |     <body>
clidownload_1  |         <h1>Telephone Book</h1>
clidownload_1  |         <hr>
clidownload_1  |         <table style="width:50%">
clidownload_1  |           <tr>
clidownload_1  |             <th>Index</th>
clidownload_1  |             <th>Name</th>
clidownload_1  |             <th>Phone</th>
clidownload_1  |             <th>Address</th>
clidownload_1  |             <th>City</th>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>0</td>
clidownload_1  |             <td>Møller</td>
clidownload_1  |             <td>&#43;45 20 86 46 44</td>
clidownload_1  |             <td>Herningvej 8</td>
clidownload_1  |             <td>4800 Nykøbing F</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>1</td>
clidownload_1  |             <td>A Egelund-Møller</td>
clidownload_1  |             <td>&#43;45 54 94 41 81</td>
clidownload_1  |             <td>Rønnebærparken 1 0011</td>
clidownload_1  |             <td>4983 Dannemare</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>2</td>
clidownload_1  |             <td>A K Møller</td>
clidownload_1  |             <td>&#43;45 75 50 75 14</td>
clidownload_1  |             <td>Bregnerødvej 75, st. 0002</td>
clidownload_1  |             <td>3460 Birkerød</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |           <tr>
clidownload_1  |             <td>3</td>
clidownload_1  |             <td>A Møller</td>
clidownload_1  |             <td>&#43;45 97 95 20 01</td>
clidownload_1  |             <td>Dalstræde 11 Heltborg</td>
clidownload_1  |             <td>7760 Hurup Thy</td>
clidownload_1  |           </tr>
clidownload_1  |
clidownload_1  |         </table>
clidownload_1  |         <p></p>
clidownload_1  |         Data taken from <a href="https://www.krak.dk/person/resultat/møller">Krak.dk</a>
clidownload_1  |     </body>
clidownload_1  | </html>
clidownload_1  |    0  45533      0 --:--:-- --:--:-- --:--:-- 45533
example_clidownload_1 exited with code 0

Killing example_dbserver_1     ... done
Killing example_webserver_1    ... done
Gracefully stopping... (press Ctrl+C again to force)
```

### Cleaning up

```bash
$ docker ps -a
CONTAINER ID        IMAGE                COMMAND                  CREATED             STATUS                       PORTS               NAMES
01a0a11d00d3        appropriate/curl     "sh -c 'sleep 5 &&..."   9 minutes ago       Exited (0) 9 minutes ago                         03containersandvms_clidownload_1
ef4617bdc0d8        helgecph/webserver   "/bin/sh -c ./telb..."   9 minutes ago       Exited (137) 6 seconds ago                       03containersandvms_webserver_1
113c782030c4        helgecph/dbserver    "docker-entrypoint..."   9 minutes ago       Exited (0) 5 seconds ago                         03containersandvms_dbserver_1
```

```bash
$ docker-compose rm
```
----------


# An example with Vagrant

See the [Vagrantfile](Vagrantfile) and start it with `vagrant up`. It will create two VMs. one with the [webserver](main.go) and the other one with the [database](db_setup.js).
