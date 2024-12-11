#!/bin/sh

kubectl port-forward -n temporal deployment/temporal-web 8233:8080
