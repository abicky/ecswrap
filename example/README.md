# ecswrap example

This example task has the following containers:

* logger container which post events to fluentd container
* fluentd container where `ecswrap` runs with `ECSWRAP_LINKED_CONTAINERS=logger`

The logger container posts events for about 5 seconds after receiving SIGTERM. The fluentd container keeps running until the logger container stops.

## Prerequisites

* awscli
* jq
* docker

## Run an example ECS task

```
./run_example_ecs_task.sh
```

The following steps will be executed by the above command:

1. Create `ecswrap-example/logger` ECR repository and `ecswrap-example/fluentd` ECR repository if they doesn't exist
2. Build docker images for the ECS task
3. Push the docker images to the ECR repositories
4. Register an ECS task definition
5. Run the task

You must specify some AWS environment variables if you receive a permission error.
cf. https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html


## Example output

### logger

```
$ docker logs $(docker ps -qaf name=logger | head -1)
E, [2019-02-17T20:45:15.774585 #1] ERROR -- : Failed to connect fluentd: Connection refused - connect(2) for "fluentd" port 24224
E, [2019-02-17T20:45:15.774639 #1] ERROR -- : Connection will be retried.
I, [2019-02-17T20:45:15.774677 #1]  INFO -- : Post an event (cnt: 1)
E, [2019-02-17T20:45:15.774774 #1] ERROR -- : Failed to post an event (cnt: 1)
I, [2019-02-17T20:45:16.775866 #1]  INFO -- : Post an event (cnt: 2)
E, [2019-02-17T20:45:16.776188 #1] ERROR -- : Failed to post an event (cnt: 2)
I, [2019-02-17T20:45:17.777281 #1]  INFO -- : Post an event (cnt: 3)
E, [2019-02-17T20:45:17.777609 #1] ERROR -- : Failed to post an event (cnt: 3)
I, [2019-02-17T20:45:18.778712 #1]  INFO -- : Post an event (cnt: 4)
E, [2019-02-17T20:45:18.778831 #1] ERROR -- : Failed to post an event (cnt: 4)
I, [2019-02-17T20:45:19.779915 #1]  INFO -- : Post an event (cnt: 5)
E, [2019-02-17T20:45:19.780212 #1] ERROR -- : Failed to post an event (cnt: 5)
I, [2019-02-17T20:45:20.781393 #1]  INFO -- : Post an event (cnt: 6)
E, [2019-02-17T20:45:20.781618 #1] ERROR -- : Failed to post an event (cnt: 6)
I, [2019-02-17T20:45:21.782706 #1]  INFO -- : Post an event (cnt: 7)
I, [2019-02-17T20:45:22.784125 #1]  INFO -- : Post an event (cnt: 8)
I, [2019-02-17T20:45:23.785450 #1]  INFO -- : Post an event (cnt: 9)
I, [2019-02-17T20:45:24.785731 #1]  INFO -- : Post an event (cnt: 10)
I, [2019-02-17T20:45:25.787034 #1]  INFO -- : Post an event (cnt: 11)
I, [2019-02-17T20:45:26.788373 #1]  INFO -- : Post an event (cnt: 12)
I, [2019-02-17T20:45:27.789637 #1]  INFO -- : Post an event (cnt: 13)
I, [2019-02-17T20:45:28.790946 #1]  INFO -- : Post an event (cnt: 14)
I, [2019-02-17T20:45:29.792244 #1]  INFO -- : Post an event (cnt: 15)
I, [2019-02-17T20:45:30.793552 #1]  INFO -- : Post an event (cnt: 16)
I, [2019-02-17T20:45:31.794871 #1]  INFO -- : Post an event (cnt: 17)
I, [2019-02-17T20:45:32.796184 #1]  INFO -- : Post an event (cnt: 18)
I, [2019-02-17T20:45:33.797460 #1]  INFO -- : Post an event (cnt: 19)
I, [2019-02-17T20:45:34.798740 #1]  INFO -- : Post an event (cnt: 20)
I, [2019-02-17T20:45:35.800027 #1]  INFO -- : Post an event (cnt: 21)
I, [2019-02-17T20:45:36.433867 #1]  INFO -- : Stopping...
I, [2019-02-17T20:45:36.801125 #1]  INFO -- : Post an event (cnt: 22)
I, [2019-02-17T20:45:37.802699 #1]  INFO -- : Post an event (cnt: 23)
I, [2019-02-17T20:45:38.804005 #1]  INFO -- : Post an event (cnt: 24)
I, [2019-02-17T20:45:39.805615 #1]  INFO -- : Post an event (cnt: 25)
I, [2019-02-17T20:45:40.807301 #1]  INFO -- : Post an event (cnt: 26)
I, [2019-02-17T20:45:41.808586 #1]  INFO -- : Stopped
```

### fluentd

```
$ docker logs $(docker ps -qaf name=fluentd | head -1)
ecswrap 2019/02/17 20:45:08 [DEBUG] opts: {Timeout:10 LinkedContainers:[logger] Verbosity:[true]}, args: [/bin/entrypoint.sh fluentd], PID: 6
ecswrap 2019/02/17 20:45:08 [DEBUG] Succeeded to start command [/bin/entrypoint.sh fluentd] with PID 10
2019-02-17 20:45:16 +0000 [info]: parsing config file is succeeded path="/fluentd/etc/fluent.conf"
2019-02-17 20:45:16 +0000 [warn]: [output_docker1] 'time_format' specified without 'time_key', will be ignored
2019-02-17 20:45:16 +0000 [warn]: [output1] 'time_format' specified without 'time_key', will be ignored
2019-02-17 20:45:16 +0000 [info]: using configuration file: <ROOT>
(snip)
</ROOT>
2019-02-17 20:45:16 +0000 [info]: starting fluentd-1.3.3 pid=10 ruby="2.5.2"
2019-02-17 20:45:16 +0000 [info]: spawn command to main:  cmdline=["/usr/bin/ruby", "-Eascii-8bit:ascii-8bit", "/usr/bin/fluentd", "-c", "/fluentd/etc/fluent.conf", "-p", "/fluentd/plugins", "--under-supervisor"]
2019-02-17 20:45:20 +0000 [info]: gem 'fluentd' version '1.3.3'
2019-02-17 20:45:20 +0000 [info]: adding match in @mainstream pattern="docker.**" type="file"
2019-02-17 20:45:20 +0000 [warn]: #0 [output_docker1] 'time_format' specified without 'time_key', will be ignored
2019-02-17 20:45:20 +0000 [info]: adding match in @mainstream pattern="**" type="file"
2019-02-17 20:45:20 +0000 [warn]: #0 [output1] 'time_format' specified without 'time_key', will be ignored
2019-02-17 20:45:20 +0000 [info]: adding filter pattern="**" type="stdout"
2019-02-17 20:45:20 +0000 [info]: adding source type="forward"
2019-02-17 20:45:20 +0000 [info]: #0 starting fluentd worker pid=21 ppid=10 worker=0
2019-02-17 20:45:20 +0000 [info]: #0 [input1] listening port port=24224 bind="0.0.0.0"
2019-02-17 20:45:20 +0000 [info]: #0 fluentd worker is now running worker=0
2019-02-17 20:45:20.371300258 +0000 fluent.info: {"worker":0,"message":"fluentd worker is now running worker=0"}
2019-02-17 20:45:20 +0000 [warn]: #0 no patterns matched tag="fluent.info"
ecswrap 2019/02/17 20:45:37 [DEBUG] container: &{Name:logger DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:38 [DEBUG] container: &{Name:fluentd DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:38 [DEBUG] container: &{Name:logger DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:39 [DEBUG] container: &{Name:logger DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:40 [DEBUG] container: &{Name:fluentd DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:40 [DEBUG] container: &{Name:logger DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:41 [DEBUG] container: &{Name:fluentd DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:41 [DEBUG] container: &{Name:logger DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:42 [DEBUG] container: &{Name:fluentd DesiredStatus:STOPPED KnownStatus:RUNNING}
ecswrap 2019/02/17 20:45:42 [DEBUG] container: &{Name:logger DesiredStatus:STOPPED KnownStatus:STOPPED}
ecswrap 2019/02/17 20:45:42 [DEBUG] All linked containers have stopped
ecswrap 2019/02/17 20:45:42 [DEBUG] Send signal 15 to child process
2019-02-17 20:45:42 +0000 [info]: Received graceful stop
2019-02-17 20:45:43 +0000 [info]: #0 fluentd worker is now stopping worker=0
2019-02-17 20:45:43.180423984 +0000 fluent.info: {"worker":0,"message":"fluentd worker is now stopping worker=0"}
2019-02-17 20:45:43 +0000 [warn]: #0 no patterns matched tag="fluent.info"
2019-02-17 20:45:43 +0000 [info]: #0 shutting down fluentd worker worker=0
2019-02-17 20:45:43 +0000 [info]: #0 shutting down input plugin type=:forward plugin_id="input1"
2019-02-17 20:45:43 +0000 [info]: #0 shutting down output plugin type=:file plugin_id="output1"
2019-02-17 20:45:43 +0000 [info]: #0 shutting down output plugin type=:file plugin_id="output_docker1"
2019-02-17 20:45:43 +0000 [info]: #0 shutting down filter plugin type=:stdout plugin_id="object:2ae6812b310c"
2019-02-17 20:45:43 +0000 [info]: Worker 0 finished with status 0
ecswrap 2019/02/17 20:45:43 [DEBUG] child process exited normally
$ docker cp $(docker ps -qaf name=fluentd | head -1):/fluentd/log/data.log .
$ docker cp $(docker ps -qaf name=fluentd | head -1):$(readlink data.log) .
$ cat data.*.log
2019-02-17T20:45:15+00:00               {"cnt":1}
2019-02-17T20:45:16+00:00               {"cnt":2}
2019-02-17T20:45:17+00:00               {"cnt":3}
2019-02-17T20:45:18+00:00               {"cnt":4}
2019-02-17T20:45:19+00:00               {"cnt":5}
2019-02-17T20:45:20+00:00               {"cnt":6}
2019-02-17T20:45:21+00:00               {"cnt":7}
2019-02-17T20:45:22+00:00               {"cnt":8}
2019-02-17T20:45:23+00:00               {"cnt":9}
2019-02-17T20:45:24+00:00               {"cnt":10}
2019-02-17T20:45:25+00:00               {"cnt":11}
2019-02-17T20:45:26+00:00               {"cnt":12}
2019-02-17T20:45:27+00:00               {"cnt":13}
2019-02-17T20:45:28+00:00               {"cnt":14}
2019-02-17T20:45:29+00:00               {"cnt":15}
2019-02-17T20:45:30+00:00               {"cnt":16}
2019-02-17T20:45:31+00:00               {"cnt":17}
2019-02-17T20:45:32+00:00               {"cnt":18}
2019-02-17T20:45:33+00:00               {"cnt":19}
2019-02-17T20:45:34+00:00               {"cnt":20}
2019-02-17T20:45:35+00:00               {"cnt":21}
2019-02-17T20:45:36+00:00               {"cnt":22}
2019-02-17T20:45:37+00:00               {"cnt":23}
2019-02-17T20:45:38+00:00               {"cnt":24}
2019-02-17T20:45:39+00:00               {"cnt":25}
2019-02-17T20:45:40+00:00               {"cnt":26}
```
