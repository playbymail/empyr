#!/bin/bash
############################################################################
#
siteRoot="a02"
[ -f "testdata/${siteRoot}/index.html" ] || {
  echo "error: missing site root"
  echo "       siteRoot == 'testdata/${siteRoot}'"
  exit 2
}

############################################################################
# create a tarball of the site that doesn't include the mac junk
tarball="site-a02.tgz"
cd "testdata" || exit 2
echo " info: creating '${tarball}..."
# tar -cz --no-xattrs --no-mac-metadata -f "${tarball}" a02 || exit 2
tar -cz --no-xattrs --no-mac-metadata --disable-copyfile --format=ustar -f "${tarball}" a02 || exit 2
echo " info: created  '${tarball}'"

############################################################################
# push the tarball to the server
echo " info: pushing '${tarball}' to server..."
scp "${tarball}" epimethean:/tmp || exit 2
echo " info: pushed  '${tarball}' to server"

############################################################################
# extract the tarball on the server
echo " info: extracting '${tarball}' on server to /var/www/a02..."
ssh epimethean "mkdir -p /var/www/a02 && tar -xzf /tmp/${tarball} -C /var/www && chown -R epimethean:epimethean /var/www/a02 && rm /tmp/${tarball}" || exit 2
echo " info: extracted  '${tarball}' on server and updated ownership"

############################################################################
#
exit 0
