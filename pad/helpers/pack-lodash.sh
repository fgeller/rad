#!/bin/bash

PACK_NAME=lodash
SOURCE=https://lodash.com/docs
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
     -D lodash.com \
     -m \
     -l1 \
     --restrict-file-names=windows \
     -p \
     -P$DOWNLOAD_DIR \
     $SOURCE

read -r -d '' PACK_DESC << EOM
Copyright 2012-2015 The Dojo Foundation <http://dojofoundation.org/>.<br />
More license information available <a href="https://github.com/lodash/lodash/blob/master/LICENSE">here</a>.<br />
lodash is available here: <a href="https://lodash.com">here</a>.<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer lodash -name $PACK_NAME -desc "$PACK_DESC" -version `date --iso` -source $DOWNLOAD_DIR -dest /tmp/packs
