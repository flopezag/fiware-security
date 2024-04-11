# Security Scan of FIWARE Catalogue components

This program has been developed to facilitate the Security Scan of the FIWARE Catalogue components and generate a report to facilitate the the resolution of identified issues on them.

The information of the components to be analyzed is maintained in the file [enablers.json](./config/enablers.json).

## Go installation

To install the Go language, you can follow the instructions detailed in the [Go Installation instructions](https://go.dev/doc/install). The following are the steps for Linux installation

1. Remove any previous Go installation by deleting the /usr/local/go folder (if it exists), then extract the archive you just downloaded into /usr/local, creating a fresh Go tree in /usr/local/go:

    ```bash
    $ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
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

It will generate the `scan` program that we will use in order t generate the summary of security vulnerabilities of our code.

## Run

To execute the scan, just specify the option of `check` together with the Enabler that we wanted to analyse. The list of available enablers can be found in the [enablers.json](./config/enablers.json) file. The command should be the following for Keyrock enabler:

```bash
scan check Keyrock
```

It will generate a file in the `Anchore`folder with the result of the Security Scan Analysis with details of the Date and Time of this scan (e.g., Keyrock_idm_20240411_1254_anchore.json) in JSON format.

Furthermore, we can use a python data analysis script to summarize the data and visualize the histogram of the different vulnerabilities found in the scan. See [helpers/README.md](../helpers/README.md) for more details.

## License

These scripts are licensed under Apache License 2.0.