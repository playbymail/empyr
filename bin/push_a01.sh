#!/bin/bash
############################################################################
#
gameCode="a01"
############################################################################
#
[ -d build -a -f "bin/push_${gameCode}.sh" ] || {
  echo error: must run from the root of the repository
  exit 2
}

############################################################################
#
echo " info: removing old builds..."
rm -rf "build/${gameCode}" "build/empyr-${gameCode}.tgz"

############################################################################
#
echo " info: creating build directory..."
mkdir -p "build/${gameCode}"/{bin,data,public,templates} || exit 2

############################################################################
#
echo " info: building executable..."
GOOS=linux GOARCH=amd64 go build -o "build/${gameCode}/bin/empyr" || exit 2
echo " info: copying public files..."
cp -R public/* "build/${gameCode}/public/" || exit 2
echo " info: copying template files..."
cp -R app/templates/* "build/${gameCode}/templates/" || exit 2

############################################################################
# create a compressed tarball of the build folder, without any Mac junk.
echo " info: creating tarball..."
cd "build/${gameCode}" || exit 2
tar -cz --no-xattrs --no-mac-metadata --disable-copyfile --format=ustar -f "../empyr-${gameCode}.tgz" bin public templates || exit 2
cd ../..

############################################################################
# push the install script and tarball to our production server
echo " info: pushing installation script..."
scp "bin/install_${gameCode}.sh" epimethean:"/var/www/${gameCode}/build/install.sh" || exit 2
echo " info: pushing tarball..."
scp "build/empyr-${gameCode}.tgz" epimethean:"/var/www/${gameCode}/build/" || exit 2

############################################################################
# execute the installation script
echo " info: executing the installation script..."
ssh epimethean "/var/www/${gameCode}/build/install.sh" || {
  echo "error: installation script failed"
  exit 2
}

############################################################################
# helpful message about the service if it is installed
if [ -f "/etc/systemd/system/empyr-${gameCode}.service" ]; then
  echo " info: if this succeeded, you should restart the services"
  echo "       ssh epimethean systemctl restart empyr-${gameCode}.service"
  echo "       ssh epimethean systemctl status  empyr-${gameCode}.service"
  echo "       ssh epimethean journalctl -f -u  empyr-${gameCode}.service"
fi

############################################################################
#
exit 0
