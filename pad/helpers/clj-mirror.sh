#!/bin/bash
set -e

PFX=clojure-stdlib
PKG=https://clojuredocs.org/core-library/vars

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
     -D clojuredocs.org \
     -m \
     -P$PFX \
     $PKG

# escape entities: <! to &lt;!
find . -name ".html" -exec sed -iE 's:<(![^dD]):\&lt;\1:g' {} \;
# cat x.html | sed -E 's:<(![^dD]):\&lt;\1:g' > '<!!.html'
