#!/bin/bash

read -r -d '' PACK_DESC << EOM
Documentation generated from sources available <a href="https://github.com/mpollmeier/gremlin-scala">here</a>.<br />
Gremlin-Scala is Open Source and available under the Apache 2 License.<br/>
Gremlin-Scala is availabe <a href="https://github.com/mpollmeier/gremlin-scala">here</a>.<br />
The Scala distribution is released under the <a href="http://www.scala-lang.org/license.html">3-clause BSD license</a>.<br />
Scala is available for download <a href="http://www.scala-lang.org/download/">here</a>.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

PACK_NAME=scala-gremlin
TARGET=d/gremlin-scala/scaladocs
mkdir -p $TARGET
find $TARGET -name "*.htmlE*" -exec rm -fv {} \;
pushd $TARGET
find ../*/src/main/scala -name "*.scala" | xargs scaladoc
popd

../pad -indexer scala -name $PACK_NAME -desc "$PACK_DESC" -version 3.0.1-incubating4 -source $TARGET -dest /tmp/packs
