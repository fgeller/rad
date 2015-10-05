#!/bin/bash
set -e

PFX=go-stdlib
#PKG=http://godoc.org/github.com/gorilla/websocket
PKG=https://godoc.org/-/go

# --wait: don't spam the host
# -E: convert extensions, adds .html/.css
# -k: convert links to relative ones
# -H: allow spanning hosts
# -D: include the following list of domains
# -m: mirror the given site
# -l: limit the recursion level
# -P: local directory prefix for downloaded files
wget --wait=0.5 \
     -E \
     -k \
     -H \
     -D'godoc.org,maxcdn.bootstrapcdn.com,ajax.googleapis.com' \
     -m \
     -l1 \
     -P$PFX \
     $PKG

# disable google analytics
find . -name "*.html" -exec sed -i -e 's:type="text/javascript":type="text/disabled":g' {} \;

# don't show page header and footer
echo ".navbar-default, #x-footer { display: none; }" >> `find $PFX -name "site.css*"`
echo "#x-pkginfo { visibility: hidden; }" >> `find $PFX -name "site.css*"`

# for some reason -X / exclude does not work
rm -v $PFX/godoc.org/-/index.html
