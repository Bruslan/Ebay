
# go gets:
go get github.com/gocql/gocql
go get golang.org/x/net/http2

- setting up single cassandra node ubuntu:
	https://www.digitalocean.com/community/tutorials/how-to-install-cassandra-and-run-a-single-node-cluster-on-ubuntu-14-04)
- cassandra install directories:
	https://docs.datastax.com/en/cassandra/latest/cassandra/install/referenceInstallLocatePkg.html
- start cassandra:
	sudo service cassandra stop
- stop cassandra:
	sudo service cassandra start
- connect to cass:
	cqlsh -u cassandra -p cassandra

0. create cassandra user and password and set cassandra superuser=false
	follow: https://docs.datastax.com/en/cassandra/latest/cassandra/configuration/secureConfigNativeAuth.html

1. call db_setup.txt using cqlsh source command
	(https://docs.datastax.com/en/cql/3.1/cql/cql_reference/source_r.html?hl=source)
	to ecxecute cassandra table setup: `/usr/bin/cqlsh -f /home/jansen/go/src/github.com/ianzn-private/data/db_setup.txt`


java Version festlegen

export JAVA_HOME=`/usr/libexec/java_home -v 1.8`



## Testing Files:
- switch into directory and execute `go test`


