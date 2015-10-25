#!/bin/bash

read -r -d '' PACK_DESC << EOM
Documentation generated from sources available <a href="https://github.com/apache/kafka">here</a>.<br />
Apache Kafka is Open Source and available under the Apache 2 License.<br/>
Apache Kafka is availabe <a href="http://kafka.apache.org/downloads.html">here</a>.<br />
The Scala distribution is released under the <a href="http://www.scala-lang.org/license.html">3-clause BSD license</a>.<br />
Scala is available for download <a href="http://www.scala-lang.org/download/">here</a>.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

PACK_NAME=kafka
TARGET=d/$PACK_NAME/scaladocs
mkdir -p $TARGET
find $TARGET -name "*.htmlE*" -exec rm -fv {} \;
pushd $TARGET
find ../core/src/main/scala -name "*.scala" | xargs scaladoc
popd

../pad -indexer scala -name scala-$PACK_NAME -desc "$PACK_DESC" -version 0.8.2.2 -source $TARGET -dest /tmp/packs
