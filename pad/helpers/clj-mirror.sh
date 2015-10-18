#!/bin/bash
PFX=clojure-stdlib
PKG=https://clojuredocs.org/core-library/vars

# --random-wait
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
     -e robots=off \
     -k \
     -H \
     -D clojuredocs.org,gravatar.com,avatars.githubusercontent.com \
     -m \
     -P$PFX \
     $PKG

echo 'header.navbar { visibility: hidden !important; display: none !important; }' >> `find $PFX -name "app.css*"`
echo 'div.desktop-side-nav { visibility: hidden !important; display: none !important; }' >>  `find $PFX -name "app.css*"`
