#!/bin/bash

read -r -d '' PACK_DESC << EOM
Documentation generated from sources available <a href="https://github.com/spray/spray">here</a>.<br />
Spray is Open Source and available under the Apache 2 License.<br/>
Spray is availabe <a href="http://spray.io">here</a>.<br />
The Scala distribution is released under the <a href="http://www.scala-lang.org/license.html">3-clause BSD license</a>.<br />
Scala is available for download <a href="http://www.scala-lang.org/download/">here</a>.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

PACK_NAME=spray
TARGET=d/$PACK_NAME/scaladocs
mkdir -p $TARGET
find $TARGET -name "*.htmlE*" -exec rm -fv {} \;
pushd $TARGET
find ../spray-*/src/main/scala -name "*.scala" | xargs scaladoc
popd

../pad -indexer scala -name scala-$PACK_NAME -desc "$PACK_DESC" -version 1.3.3 -source $TARGET -dest /tmp/packs
