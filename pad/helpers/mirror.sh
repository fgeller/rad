#!/bin/bash
set -e

PFX=godoc

# --wait: don't spam the host
# -E: convert extensions, adds .html/.css
# -k: convert links to relative ones
# -H: allow spanning hosts
# -D: include the following list of domains
# -D: exclude the listed directories
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
     http://godoc.org/-/go

# don't show page header and footer
echo ".navbar-default, #x-footer, #x-pkginfo p { display: none; }" >> `find $PFX -name "site.css*"`

rm -v $PFX/godoc.org/-/index.html
