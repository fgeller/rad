#!/bin/bash

read -r -d '' PACK_DESC << EOM
Documentation generated from sources available <a href="https://github.com/milessabin/shapeless">here</a>.<br />
Shapeless is Open Source and available under the Apache 2 License.<br/>
Shapeless is availabe <a href="https://github.com/milessabin/shapeless">here</a>.<br />
The Scala distribution is released under the <a href="http://www.scala-lang.org/license.html">3-clause BSD license</a>.<br />
Scala is available for download <a href="http://www.scala-lang.org/download/">here</a>.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

PACK_NAME=shapeless
TARGET=d/shapeless/scaladocs
mkdir -p $TARGET
find $TARGET -name "*.htmlE*" -exec rm -fv {} \;
pushd $TARGET
find ../core/src/main/scala -name "*.scala" | xargs scaladoc
popd

../pad -indexer scala -name scala-$PACK_NAME -desc "$PACK_DESC" -version 2.2.5 -source $TARGET -dest /tmp/packs
