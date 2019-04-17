#!/bin/bash

#  return ip or hostname of external LB

if [ $# -lt 2 ]; then
    echo "Usage: $0 <k8s namespace> <k8s ingress-gateway>"
    exit
fi

gw=$(kubectl -n "$1" get svc "$2" \
    -o jsonpath="{.status.loadBalancer.ingress[0].ip}")

if [ -z $gw ]; then
    gw=$(kubectl -n "$1" get svc "$2" \
        -o jsonpath="{.status.loadBalancer.ingress[0].hostname}")
fi

echo $gw