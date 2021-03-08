#!/bin/bash
kubectl create ns go-admin
kubectl create configmap settings-admin --from-file=../../config/settings.yml -n go-admin
