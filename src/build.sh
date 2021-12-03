#!/bin/bash
version=1.0
#if [ $# -eq 0 ] 
#then
#        echo "Please input version, like \"./release.sh 0.60\""
#        exit
#fi
rm -f release/portscan_*$version.tgz
echo "Build ReleaseFile for version $version"

#export GOPATH=`pwd`

echo "build linux_amd64"
export GOOS=linux GOARCH=amd64 
go build -ldflags="-w -s"
tar zcvf portscan_linux_amd64_$version.tgz portscan
rm -f portscan  portscan.exe 

echo "build linux_386"
export GOOS=linux GOARCH=386 
go build -ldflags="-w -s"
tar zcvf portscan_linux_386_$version.tgz portscan
rm -f portscan  portscan.exe 

echo "build mac_x64"
export GOOS=darwin GOARCH=amd64 
go build -ldflags="-w -s"
tar zcvf portscan_mac_amd64_$version.tgz portscan
rm -f portscan  portscan.exe 

echo "build mac_arm64"
export GOOS=darwin GOARCH=arm64 
go build -ldflags="-w -s"
tar zcvf portscan_mac_arm64_$version.tgz portscan
rm -f portscan  portscan.exe 

echo "build win32"
export GOOS=windows GOARCH=386 
go build -ldflags="-w -s"
tar zcvf portscan_win32_$version.tgz portscan.exe 
rm -f portscan  portscan.exe 

echo "build win64"
export GOOS=windows GOARCH=amd64 
go build -ldflags="-w -s"
tar zcvf portscan_win64_$version.tgz portscan.exe 
rm -f portscan  portscan.exe 

echo "Build Over"

mkdir release
mv *.tgz release
ls -l release/portscan_*$version.tgz