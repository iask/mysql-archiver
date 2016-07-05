# build
go build -ldflags "-w" -o bin/archiver archiver.go 

# init env
  1 login your all mysql instance and execute grants

    GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REPLICATION SLAVE ON *.* TO 'ptarchiver'@'your archiver admin IP' IDENTIFIED BY 'xxx';

  2 import admin schema

    mysql -uroot -p < archiver.sql


# run api
./bin/archiver --conf=conf/app.conf

## web doc
  # list all scheduleds
  GET : http://127.0.0.1:9091/archiver/schds

  # list all jobs
  GET: http://127.0.0.1:9091/archiver/jobs

  # list all cron
  GET: http://127.0.0.1:9091/archiver/crons


## api doc
  # list all scheduleds
  GET : http://127.0.0.1:9090/v1/schds

  # get a scheduleds
  GET : http://127.0.0.1:9090/v1/schds/123

  # dryrun a scheduled
  GET : http://127.0.0.1:9090/v1/schds/123/1

  # run a scheduled
  GET : http://127.0.0.1:9090/v1/schds/123/2
 
  # add a scheduled
  POST: http://127.0.0.1:9090/v1/schds

  # update a scheduleds
  POST : http://127.0.0.1:9090/v1/schds/123

  # delete a scheduled
  DELETE: http://127.0.0.1:9090/v1/schds/123

  # list all jobs
  GET: http://127.0.0.1:9090/v1/jobs

  # list a jobs log
  GET: http://127.0.0.1:9090/v1/jobs/345/log

  # list all cron
  GET: http://127.0.0.1:9090/v1/crons

