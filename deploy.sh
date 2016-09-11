#!/bin/bash

echo "delete rpi_temp daemon set (ds)"
kubectl delete -f rpi_temp.yaml

echo "delete old docker tag definition"
docker rmi 10.0.0.20:5000/mesgtone/rpi_temp

echo "create new docker tag definiton"
docker tag rpi_temp 10.0.0.20:5000/mesgtone/rpi_temp

echo "push docker image in the Kubernetes based repository"
docker push 10.0.0.20:5000/mesgtone/rpi_temp

echo "create rpi_temp daemon set (ds)"
kubectl create -f rpi_temp.yaml

echo "sleep for 10 sec"
sleep 10
echo "check rpi_temp pods"
kubectl get pods -o wide

