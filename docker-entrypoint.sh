#!/bin/bash

mkdir -p /var/cache/bind/
cp -r /zone/* /var/cache/bind/
chown -R bind:bind /var/cache/bind/

if ! [ -f "$UPDATEKEY" ]
then
    rndc-confgen | sed -n '2,5p' > "$UPDATEKEY"
fi

bash init.sh

service named start
./nasa-judge-named

