# FIWARE Docker Security Scan

## Overview

Automatically scan a particular local docker image or all local docker containers 
with [Clair Vulnerability Scanner](https://github.com/coreos/clair) using 
[Clair-Scanner](https://github.com/arminc/clair-scanner) and 
[clair-local-scan](https://github.com/arminc/clair-local-scan) together with together 
with the Docker Bench for Security to check common best-practices around deploying 
FIWARE Docker containers in production. 

The tests are all automated, and are inspired by the 
[CIS Docker Community Edition Benchmark v1.1.0](https://benchmarks.cisecurity.org/tools2/docker/CIS_Docker_Community_Edition_Benchmark_v1.1.0.pdf).


## Installation

There are two ways to install and execute the code. The first one is installing locally
the configuration files and script to execute the [docker-compose](https://docs.docker.com/compose/) 
locally (see [README.md](docker/README.md)) or [Ansible](https://www.ansible.com/) to deploy 
a virtual machine inside [FIWARE Lab](https://cloud.lab.fiware.org) and preconfigure all 
the system to launch the scan automatically (see [README.md](deploy/README.md)).


## Credits

* Docker
* docker-compose
* [Clair Vulnerability Scanner](https://github.com/coreos/clair)
* [Clair-Scanner](https://github.com/arminc/clair-scanner) (release v8 is included)
* [clair-local-scan](https://github.com/arminc/clair-local-scan)
* [clair-container-scan](https://github.com/usr42/clair-container-scan)
* [Docker Bench Security](https://github.com/docker/docker-bench-security)

## License

These scripts are licensed under Apache License 2.0.