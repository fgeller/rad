#!/bin/bash

PACK_NAME=clojure-stdlib
SOURCE=https://clojuredocs.org/core-library/vars
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
     -D clojuredocs.org,gravatar.com,avatars.githubusercontent.com \
     -m \
     -P$DOWNLOAD_DIR \
     $SOURCE

echo 'header.navbar { visibility: hidden !important; display: none !important; }' >> `find $PFX -name "app.css*"`
echo 'div.desktop-side-nav { visibility: hidden !important; display: none !important; }' >>  `find $PFX -name "app.css*"`

read -r -d '' PACK_DESC << EOM
Examples and other content provided by <a href="https://clojuredocs.org">clojuredocs.org</a>.<br />
Clojure is available under the <a href="http://opensource.org/licenses/eclipse-1.0.php">Eclipse Public License 1.0</a>.<br />
More information available <a href="http://clojure.org/license">here</a>.<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer clojure -name $PACK_NAME -desc "$PACK_DESC" -version `date --iso` -source $DOWNLOAD_DIR -dest /tmp/packs
