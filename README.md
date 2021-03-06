# Biome
Biome is a tool that allows you to configure a temporary environment while running a command. This project was inspired by [awsudo](https://github.com/makethunder/awsudo) but adds additional functionality such as additional environment variables via a `.biome.yaml` configuration file.

## Installation
Download the release for your target platform from the [releases](https://github.com/jeff-roche/biome/releases/latest) page and add it to your `$PATH`

### Go install
`go install github.com/jeff-roche/biome`

## Configuration
Biome gets its configuration from a `.biome.yaml` file. It will look for the `.biome.yaml` file in the current directory (where the command is run) first. If it can't find the file there, biome will look in the current users home directory for the file.

### `.biome.yaml` format
As the extension shows, biome uses yaml for it's configuration format. Here is an example configuration which can also be seen in [example.biome.yaml](./example.biome.yaml).

```yaml
# .biome.yaml

name: my-biome # Biome name, required
environment:
    MY_USEFUL_ENV: "A value I need"
    MY_OTHER_ENV: "Another value I need"
```

### Dotenv File
`.env` files can be loaded in by specifying the `load_env` tag. Any vars specified in the `environment` section will override values set in the dotenv file specified.

```yaml
# .biome.yaml
name: my-biome
load_env: my_env_file.env # Specify the name of the file to load in
environment:
    MY_USEFUL_ENV: "A value I need"
    MY_OTHER_ENV: "Another value I need"
```

### AWS Environment
By specifying the `aws_profile` configuration value, Biome will load that [AWS Profile](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) from `~/.aws/credentials` and configure the environment variables and a session for this command.

```yaml
# .biome.yaml
name: my-biome
aws_profile: my_aws_profile # Important part
environment:
    ...
```

### Commands
Additional commands can be run using the commands setting. Any commands specified will be run as the last steps prior to running the top level command specified when running biome.

```yaml
# .biome.yaml
name: my-biome
commands:
    - kubectx my-k8s-context
environment:
    ...
```


## Usage
The most common use case is for use with scripts that need context via environment variables. The need for this tool came about for CI/CD scripts that need AWS context as well as additional environment variables that change based on certain states. This tool will allow you to configure those different states and provide that context to your scripts and pipelines.

### Via direct command
```bash
$ biome run -b my-biome ${COMMAND}
```

- `-b` is a required parameter that specifies the name of the biome you want to use
    - In this case, the name of the biome is `my-biome`
- `${COMMAND}` specifies the command you want to run (ex: `env`, `ls -al`, ``, etc.)

> **NOTE** if you want to add command line flags to the command being run, you need to preface it with `--`
    *Example*: `biome run -b my-biome -- ls -al`

### Via bash alias
A way that makes Biome a little more convenient is to alias your profiles via bash aliases and use them that way.

This configuration:
```bash
# ~/.bashrc
alias onstaging='biome run -b staging-biome'
alias onprod='biome run -b production-biome'
```

Allows the following command to be run on the command line:

```bash
$ onstaging ./bin/ci/deploy-service.sh
```

### Exporting to a dotenv file
You can export the loaded environment variables to a dotenv file with the following command:
```bash
$ biome save -b my-biome -f my.env
```

> **NOTE**: If no input file is specified with `-f`, `./.env` will be used as the file path

> **NOTE**: AWS environment variables are not currently exported

## Future Plans
- Have goreleaser create a docker image and publish to ghcr
- :white_check_mark: Custom MFA token message for AWS profiles
- :white_check_mark: Export loaded variables to a dotenv file
- :white_check_mark: Switch to [cobra](https://github.com/spf13/cobra) for the cli
- :white_check_mark: Allow CLI input to be a secret (for passwords)
- :white_check_mark: Allow inhereting from other biomes in the same file
- :white_check_mark: Allow setting an environment variable from stdin
- :white_check_mark: Implement goreleaser for binary building
    - :white_check_mark: Use semantic versioning
- :white_check_mark: Add a version command
- :white_check_mark: Accept all valid yaml file extensions
- :white_check_mark: Build a CI/CD pipeline
- :white_check_mark: Implement some tests
- :white_check_mark: Loading Environment variables from a .env file
- :white_check_mark: Encrypted environment variables via [dragoman](https://github.com/meltwater/dragoman)
- :white_check_mark: Kubernetes context setting
    - *NOTE* this is done through commands