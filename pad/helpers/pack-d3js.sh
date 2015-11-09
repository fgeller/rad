#!/bin/bash

PACK_NAME=d3js
SOURCE=https://github.com/mbostock/d3/wiki/API-Reference
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
     -D github.com \
     -m \
     -l 1 \
     -p \
     --restrict-file-names=windows \
     -P$DOWNLOAD_DIR \
     $SOURCE

# TODO
# disable http://collector.githubapp.com/github/page_view?dimensions[page]=http%3A%2F…8-478-42-0-0-0---0---9-13-10--14-14&&&dimensions[cid]=412254477.1446872816

# find $DOWNLOAD_DIR -name "*.js" -exec sed -i 's/null==window.GitHub/false/g' {} \;
find $DOWNLOAD_DIR -name "*.js" -exec sed -i 's/top!==window/false/g' {} \;

echo ''                                      >> `find $DOWNLOAD_DIR -name "github-*.css"`
echo '.header { display: none }'             >> `find $DOWNLOAD_DIR -name "github-*.css"`
echo '.pagehead { display: none }'           >> `find $DOWNLOAD_DIR -name "github-*.css"`
echo '.repository-sidebar { display: none }' >> `find $DOWNLOAD_DIR -name "github-*.css"`
echo '#wiki-rightbar { display: none }'      >> `find $DOWNLOAD_DIR -name "github-*.css"`
echo '.gh-header-meta { display: none }'     >> `find $DOWNLOAD_DIR -name "github-*.css"`

read -r -d '' PACK_DESC << EOM
d3 was created by <a href="http://bost.ocks.org/mike/">Mike Bostock</a>.<br />
More information about d3 is available <a href="https://d3js.org/">here</a>.<br />
<br />
Mirrored from <a href="$SOURCE">$SOURCE</a>.<br />
Pack created on `date`.<br />
<br />
EOM

../pad -indexer d3 -name $PACK_NAME -desc "$PACK_DESC" -version `date --iso` -source $DOWNLOAD_DIR -dest /tmp/packs
