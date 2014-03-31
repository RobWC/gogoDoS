[![Build Status](https://magnum.travis-ci.com/JNPRAutomate/gogoDoS.svg?token=Taq81d9PL7keqp96e9qu&branch=master)](https://magnum.travis-ci.com/JNPRAutomate/gogoDoS)

gogoDoS
=======

DoS testing tool

This tool helps in testing out our DoS prevention products. While this tool could support multiplatform as of today it is Linux only.

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

