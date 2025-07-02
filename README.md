# FIWARE Cybersecurity Analysis

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/flopezag/fiware-security">
    <img src="go/doc/FIWARESecurity.png" alt="Logo" width="137" height="150">
  </a>

<h3 align="center">FIWARE Cybersecurity Analysis of the FIWARE Generic Enablers</h3>

  <p align="center">
    <!--<a href="https://github.com/flopezag/fiware-security"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/flopezag/fiware-security">View Demo</a>
    ·-->
    <a href="https://github.com/flopezag/fiware-security/issues">Report Bug</a>
    ·
    <a href="https://github.com/flopezag/fiware-security/issues">Request Feature</a>
  </p>
</div>

# Security Scan of FIWARE Catalogue components

This program has been developed to facilitate the Security Scan of the FIWARE Catalogue components and generate a report to facilitate the the resolution of identified issues on them.

Automatically scan a particular local docker image or all local docker containers 
with [Clair Vulnerability Scanner](https://github.com/coreos/clair) using 
[Clair-Scanner](https://github.com/arminc/clair-scanner) and 
[clair-local-scan](https://github.com/arminc/clair-local-scan) together with together 
with the [Docker Bench for Security](https://github.com/docker/docker-bench-security) 
to check common best-practices around deploying FIWARE Docker containers in production. 

The tests are all automated, and are inspired by the 
[CIS Docker Community Edition Benchmark v1.1.0](https://benchmarks.cisecurity.org/tools2/docker/CIS_Docker_Community_Edition_Benchmark_v1.1.0.pdf).

The information of the components to be analyzed is maintained in the file [enablers.json](./config/enablers.json).

## Go installation

To install the Go language, you can follow the instructions detailed in the [Go Installation instructions](https://go.dev/doc/install). The following are the steps for Linux installation

1. Remove any previous Go installation by deleting the /usr/local/go folder (if it exists), then extract the archive you just downloaded into /usr/local, creating a fresh Go tree in /usr/local/go:

    ```bash
    $ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz
    ```

    (You may need to run the command as root or through sudo).

    Do not untar the archive into an existing /usr/local/go tree. This is known to produce broken Go installations.

2. Add /usr/local/go/bin to the PATH environment variable.
You can do this by adding the following line to your $HOME/.profile or /etc/profile (for a system-wide installation):

    ```bash
    export PATH=$PATH:/usr/local/go/bin
    ```

    Note: Changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just run the shell commands directly or execute them from the profile using a command such as source $HOME/.profile.

3. Verify that you've installed Go by opening a command prompt and typing the following command:

    ```bash
    $ go version
    ```

    Confirm that the command prints the installed version of Go.

## Update dependencies

To update the current dependencies of the project, execute the following command:

```bash
go mod tidy
```

## Compile the program

The command to generate the executable command of the parser is the following:

```bash
go build .
```

It will generate the `scan` program that we will use to generate the summary of security vulnerabilities of our code.

## Run

To execute the scan, just specify the option of `check` together with the Enabler that we wanted to analyse. The list of available enablers can be found in the [enablers.json](./config/enablers.json) file. The command should be the following for Keyrock enabler:

```bash
scan check Keyrock
```

It will generate a file in the `results`folder with the result of the Security Scan Analysis with details of the Date and Time of this scan (e.g., Keyrock_idm_20240411_1254_grype.json) in JSON format.

Furthermore, we can use a other command to summarize the data and visualize the histogram of the different vulnerabilities found in the scan.

```bash
scan visualize Keyrock
```
This provide console output with teh following content:

- Total count of vulnerabilities
- Severity test histogram
- EPSS and risk averages
- Count of EPSS > 0.9 and Risk > 90

where:

- **Severity**: String severity based on CVSS scores and indicate the significance 
of a vulnerability in levels. This balances concerns such as ease of exploitability, 
and the potential to affect confidentiality, integrity, and availability of software 
and services.

- **EPSS**: [Exploit Prediction Scoring System](https://www.first.org/epss/model) is 
a metric expressing the likelihood that a vulnerability will be exploited in the wild 
over the next 30 days (on a 0–1 scale); higher values signal a greater likelihood of 
exploitation. The table output shows the EPSS percentile, a one-way transform of the 
EPSS score showing the proportion of all scored vulnerabilities with an equal or lower 
probability. Percentiles linearize a heavily skewed distribution, making threshold 
choice (e.g. “only CVEs above the 90th percentile”) straightforward.

## License

These scripts are licensed under Apache License 2.0.