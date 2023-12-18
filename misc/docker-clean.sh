#!/bin/bash

SELECTOR=libhighctidh
CONTAINERS=$(docker image ls |grep $SELECTOR|cut -f1 -d\  );

for c in $CONTAINERS
do
    echo "Removing $c...";
    docker rmi $c;
    echo "Done removing $c";
done
