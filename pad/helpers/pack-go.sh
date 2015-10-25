#!/bin/bash

PACK_NAME=go-stdlib
SOURCE=https://godoc.org/-/go
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
     -D'godoc.org,maxcdn.bootstrapcdn.com,ajax.googleapis.com' \
     -m \
     -l1 \
     -P$DOWNLOAD_DIR \
     $SOURCE

echo "Disabling Google analytics"
find $DOWNLOAD_DIR -name "*.html" -exec sed -i -E 's:type="text/javascript":type="text/disabled":g' {} \;

echo "Hiding header and footer"
echo ".navbar-default, #x-footer { display: none; }" >> `find $DOWNLOAD_DIR -name "site.css*"`
echo "#x-pkginfo { visibility: hidden; }" >> `find $DOWNLOAD_DIR -name "site.css*"`

# for some reason -X / exclude does not work
echo "Manually deleting huge package index."
rm -fv $DOWNLOAD_DIR/godoc.org/-/index.html

read -r -d '' PACK_DESC << EOM
The contents of <a href="http://golang.org">golang.org</a> is licensed under the Creative Commons Attribution 3.0 License.<br />
More details on the page and <a href="https://developers.google.com/site-policies#restrictions">here</a>.<br />
Go itself is distributed with a <a href="https://golang.org/LICENSE">BSD-style license</a>.<br /><br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer go -name $PACK_NAME -desc "$PACK_DESC" -version `date --iso` -source $DOWNLOAD_DIR -dest /tmp/packs
