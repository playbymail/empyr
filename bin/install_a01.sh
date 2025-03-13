#!/bin/bash
############################################################################
#
gameCode="a01"
webRoot="/var/www/${gameCode}"
echo " info: webRoot      == '${webRoot}'"
############################################################################
#
if [ ! -d "${webRoot}" ]; then
  echo "error: missing web root"
  exit 2
elif [ ! -d "${webRoot}/bin" ]; then
  echo "error: missing web bin folder"
  exit 2
elif [ ! -d "${webRoot}/build" ]; then
  echo "error: missing web build folder"
  exit 2
elif [ ! -d "${webRoot}/data" ]; then
  echo "error: missing web data folder"
  exit 2
elif [ ! -d "${webRoot}/public" ]; then
  echo "error: missing web public folder"
  exit 2
elif [ ! -d "${webRoot}/templates" ]; then
  echo "error: missing web templates folder"
  exit 2
elif [ ! -f "${webRoot}/build/empyr-${gameCode}.tgz" ]; then
  echo "error: missing tarball"
  exit 2
fi
############################################################################
#
echo " info: setting def to web root..."
cd "${webRoot}"  || exit 2
############################################################################
#
if [ -f "${webRoot}/bin/empyr" ]; then
  echo " info: removing old executable..."
  rm "${webRoot}/bin/empyr" || exit 2
fi
############################################################################
#
echo " info: extracting tarball..."
tar xzf "build/empyr-${gameCode}.tgz" || exit 2
echo " info: forcing bits on executable..."
chmod 755 "${webRoot}/bin/empyr" || exit 2
echo " info: forcing root ownership on web root..."
chown -R root:root "${webRoot}/"{bin,public,templates} || exit 2
echo " info: resetting ownership on data folder..."
chown -R empyr:empyr "${webRoot}/data" || exit 2
############################################################################
#
echo " info: testing executable..."
"${webRoot}/bin/empyr" version || exit 2
############################################################################
#
echo " info: removing tarball..."
rm "build/empyr-${gameCode}.tgz" || exit 2
############################################################################
#
echo " info: installation completed successfully"
exit 0
