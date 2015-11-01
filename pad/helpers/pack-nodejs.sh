#!/bin/bash

PACK_NAME=nodejs
SOURCE=https://nodejs.org/api/
DOWNLOAD_DIR=d/$PACK_NAME

rm -vrf $DOWNLOAD_DIR
mkdir -p $DOWNLOAD_DIR

# --random-wait
# --wait: don't spam the host
# -E: convert extensions, adds .html/.css
# -k: convert links to relative ones
# -H: allow spanning hosts
# -D: include the following list of domains
# -m: mirror the given site
# -l: limit the recursion level
# -P: local directory prefix for downloaded files
wget --random-wait \
     -E \
     -e robots=off \
     -k \
     -H \
     --no-check-certificate \
     -D nodejs.org,fonts.gstatic.com,fonts.googleapis.com \
     -m \
     -l 1 \
     --restrict-file-names=windows \
     -p \
     -P$DOWNLOAD_DIR \
     $SOURCE

rm -fv $DOWNLOAD_DIR/index.html
rm -fv $DOWNLOAD_DIR/nodejs.org/api/documentation.html
rm -fv $DOWNLOAD_DIR/nodejs.org/api/all.html
echo '#column2, #toc, #column1 header { display: none; }' >> `find $DOWNLOAD_DIR -name '*style.css*'`

read -r -d '' PACK_DESC << EOM
License information for node.js is available <a href="https://github.com/nodejs/node/blob/master/LICENSE">here</a>.<br />
node.js is available for download <a href="https://nodejs.org/en/download/">here</a>.<br/>
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer nodejs -name $PACK_NAME -desc "$PACK_DESC" -version 4.2.1 -source $DOWNLOAD_DIR -dest /tmp/packs
