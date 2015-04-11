Porter: docker container manager via consul config
=================================================

## Default values:
* **docker sock**: `/var/run/docker.sock`
* **consul api**: `http://localhost:8500`

## Commands:

|Command | Param | Description |
|--------|-------|-------------|
|cleanup| | Removes images without names and containers with status Exited |
|run | serviceName | Creates and runs container via consul config |


