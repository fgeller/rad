#!/bin/bash
# PFX=go-stdlib
# PKG=https://godoc.org/-/go
PFX=gorilla-websocket
PKG=http://godoc.org/github.com/gorilla/websocket

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
     -P$PFX \
     $PKG

echo "Disabling Google analytics"
find $PFX -name "*.html" -exec sed -i -E 's:type="text/javascript":type="text/disabled":g' {} \;

echo "Hiding header and footer"
echo ".navbar-default, #x-footer { display: none; }" >> `find $PFX -name "site.css*"`
echo "#x-pkginfo { visibility: hidden; }" >> `find $PFX -name "site.css*"`

# for some reason -X / exclude does not work
echo "Manually deleting huge package index."
rm -fv $PFX/godoc.org/-/index.html
