#!/bin/bash

read -r -d '' PACK_DESC << EOM
Documentation downloaded from <a href="http://www.scala-lang.org/download/2.11.7.html">here</a>.<br />
The Scala distribution is released under the <a href="http://www.scala-lang.org/license.html">3-clause BSD license</a>.<br />
Scala is available for download <a href="http://www.scala-lang.org/download/">here</a>.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

for PACK_NAME in scala-library scala-compiler scala-reflect
do
    DOWNLOAD_DIR=d/scala-docs-2.11.7/api/$PACK_NAME
    ../pad -indexer scala -name $PACK_NAME -desc "$PACK_DESC" -version 2.11.7 -source $DOWNLOAD_DIR -dest /tmp/packs
done
