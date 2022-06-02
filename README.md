# Biome
Biome is a tool that allows you to configure a temporary environment while running a command. This project was inspired by [awsudo](https://github.com/makethunder/awsudo) but adds additional functionality such as additional environment variables via a `.biome.yaml` configuration file.

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

### AWS Environment
By specifying the `aws_profile` configuration value, Biome will load that [AWS Profile](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) from `~/.aws/credentials` and configure the environment variables and a session for this command.

```yaml
# .biome.yaml
name: my-biome
aws_profile: my_aws_profile # Important part
environment:
    ...
```


## Usage
The most common use case is for use with scripts that need context via environment variables. The need for this tool came about for CI/CD scripts that need AWS context as well as additional environment variables that change based on certain states. This tool will allow you to configure those different states and provide that context to your scripts and pipelines.

### Via direct command
```bash
$ biome -b my-biome ${COMMAND}
```

- `-b` is a required parameter that specifies the name of the biome you want to use
    - In this case, the name of the biome is `my-biome`
- `${COMMAND}` specifies the command you want to run (ex: `env`, `ls -al`, ``, etc.)

### Via bash alias
A way that makes Biome a little more convenient is to alias your profiles via bash aliases and use them that way.

This configuration:
```bash
# ~/.bashrc
alias onstaging='biome -b staging-biome'
alias onprod='biome -b production-biome'
```

Allows the following command to be run on the command line:

```bash
$ onstaging ./bin/ci/deploy-service.sh
```

## Future Plans
- :white_check_mark: Implement goreleaser for binary building
    - :white_check_mark: Use semantic versioning
- Add a version command
- Accept all valid yaml file extensions
- Build a Drone CI/CD pipeline
- Implement some tests
- Loading Environment variables from a .env file
- Encrypted environment variables via [dragoman](https://github.com/meltwater/dragoman)
- Potentially switching to [cobra](https://github.com/spf13/cobra) for the cli