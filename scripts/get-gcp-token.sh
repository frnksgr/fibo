#!/bin/bash

if token=$(gcloud config config-helper --format 'value(credential.id_token)' 2>/dev/null); then
    echo $token
else
    echo "failed to get gcp token"
fi
