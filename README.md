gogoDoS
=======

DoS testing tool

This tool helps in testing out our DoS product portfolio

Buidling:

export GOPATH=`pwd`

cd src

go build gogoDoS.go

Running gogoDoS

gogoDoS -p dns -d 8.8.8.8 -P 53 -r 30
 
  Opts:

  -p dns #Protocol to use today dns only

  -d 1.2.3.4 #Server or list of servers to query against 1.2.3.4 or 1.2.3.4,2.3.4.5

  -r 30 #Rate of queries per second

  -d 60 #Duration of how many seconds to run

