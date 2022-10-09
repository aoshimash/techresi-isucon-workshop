#!/bin/bash

# app1
multipass launch --bridged --name private-isu-app-1 --cpus 1 --disk 16G --mem 2G --cloud-init app.cfg 20.04 &
# app2
multipass launch --bridged --name private-isu-app-2 --cpus 1 --disk 16G --mem 2G --cloud-init app.cfg 20.04 &
# app3
multipass launch --bridged --name private-isu-app-3 --cpus 1 --disk 16G --mem 2G --cloud-init app.cfg 20.04 &
# bench
multipass launch --bridged --name private-isu-bench --cpus 2 --disk 16G --mem 4G --cloud-init benchmarker.cfg 20.04
