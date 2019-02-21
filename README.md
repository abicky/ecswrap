# ecswrap

`ecswrap` is a wrappper program for ECS containers and resolves termination order (cf. [amazon-ecs-agent#474](https://github.com/aws/amazon-ecs-agent/issues/474)) by specifying linked containers as a command line option or an environment variable.

## Usage

```
Usage:
  ecswrap [OPTIONS] -- COMMAND [ARGS]

Application Options:
      --stop-wait-timeout=       Maximum time duration in seconds to wait from when the process receives SIGTERM before sending SIGTERM to the child. This value should be less than
                                 ECS_CONTAINER_STOP_TIMEOUT. (default: 10) [$ECSWRAP_STOP_WAIT_TIMEOUT]
      --linked-container=        container names linked with the container where this program is running. [$ECSWRAP_LINKED_CONTAINERS]
      --signal-forwarding-delay= Delay seconds until forwarding a signal, which is SIGTERM, SIGQUIT or SIGINT,to child processes. (default: 0) [$ECSWRAP_SIGNAL_FORWARDING_DELAY]
  -v, --verbose                  Verbosity

Help Options:
  -h, --help                     Show this help message
```

## Example

See [example](example).

## Author

Takeshi Arabiki (abicky)
