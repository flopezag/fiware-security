# FIWARE Docker Security Scan

## Overview

Automatically scan a particular local docker image or all local docker containers 
with [Clair Vulnerability Scanner](https://github.com/coreos/clair) using 
[Clair-Scanner](https://github.com/arminc/clair-scanner) and 
[clair-local-scan](https://github.com/arminc/clair-local-scan) together with together 
with the [Docker Bench for Security](https://github.com/docker/docker-bench-security) 
to check common best-practices around deploying FIWARE Docker containers in production. 

The tests are all automated, and are inspired by the 
[CIS Docker Community Edition Benchmark v1.1.0](https://benchmarks.cisecurity.org/tools2/docker/CIS_Docker_Community_Edition_Benchmark_v1.1.0.pdf).


## Installation

There are two ways to install and execute the code. The first one is installing locally
the configuration files and script to execute the [docker-compose](https://docs.docker.com/compose/) 
locally (see [README.md](docker/README.md)) or [Ansible](https://www.ansible.com/) to deploy 
a virtual machine inside [FIWARE Lab](https://cloud.lab.fiware.org) and preconfigure all 
the system to launch the scan automatically (see [README.md](deploy/README.md)).

## Adding the Anchore Scan GitHub Action Workflow to a repository directly

Anchore provide a [GitHub Action](https://github.com/anchore/scan-action) for Vulnerability Scanning.
Two sample GitHub Action Workflows have been added to this repository.

For example, to enable an Anchore Scan of a Docker image based on **node-slim**:

-  Copy the [anchore-node-slim.yaml](https://github.com/flopezag/fiware-security/blob/master/.github/workflows/anchore-node-slim.yml) file to `.github/workflows/anchore-node-slim.yml
-  Amend the Dockerfile [context location](https://github.com/flopezag/fiware-security/blob/master/.github/workflows/anchore-node-slim.yml#L34) if necessary - the example assumes a folder called `docker` is used.
-  After committing and pushing the file, run the new GitHub Action [manually](https://docs.github.com/en/actions/managing-workflow-runs/manually-running-a-workflow)

A security report will be displayed on
`https://github.com/<Owner>/<Repository>/security/code-scanning?query=is%3Aopen+branch%3Amaster+severity%3Aerror`


![](./img/alerts.png)

Like any GitHub Action Workflow, the creation of additional Docker images to scan
can also be added to a repository and creation can be arbitrarily more complex.
A second example file shows how to build an [alternative base image](https://kuberty.io/blog/best-os-for-docker/)
using `--build-arg` parameters on the command line to create a container based on [Red Hat UBI (Universal Base Image) 8](https://developers.redhat.com/articles/2021/11/08/optimize-nodejs-images-ubi-8-nodejs-minimal-image). To scan this alternate image, just copy
over [anchore-ubi.yaml](https://github.com/flopezag/fiware-security/blob/master/.github/workflows/anchore-ubi.yml) to `.github/workflows/anchore-ubi.yml`

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