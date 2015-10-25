#!/bin/bash

read -r -d '' PACK_DESC << EOM
Documentation downloaded from <a href="http://www.oracle.com/technetwork/java/javase/documentation/jdk8-doc-downloads-2133158.html">here</a>.<br />
The Java SE Development Kit 8 Documentation License Agreement can be found <a href="http://www.oracle.com/technetwork/java/javase/overview/javase8speclicense-2158700.html">here</a>.<br />
Java SE is available for download <a href="http://www.oracle.com/technetwork/java/javase/downloads/index.html">here</a>.<br/>
<br/>
Pack created on `date`.<br />
<br />
EOM

PACK_NAME=java
DOWNLOAD_DIR=d/jdk-8u60-docs/api
../pad -indexer java -name $PACK_NAME -desc "$PACK_DESC" -version 8u60 -source $DOWNLOAD_DIR -dest /tmp/packs
