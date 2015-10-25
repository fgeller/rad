#!/bin/bash

read -r -d '' PACK_DESC << EOM
Documentation downloaded from <a href="http://akka.io/downloads/">here</a>.<br />
Akka is available for download from <a href="http://akka.io/downloads/">here</a>.<br />
Akka is Open Source and available under the Apache 2 License.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

PACK_NAME=akka
DOWNLOAD_DIR=d/akka-2.4.0/doc/akka/api
find $DOWNLOAD_DIR -name "*.htmlE" -exec rm -fv {} \;
../pad -indexer scala -name $PACK_NAME -desc "$PACK_DESC" -version 2.4.0 -source $DOWNLOAD_DIR -dest /tmp/packs
