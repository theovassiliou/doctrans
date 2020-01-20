# Content of this directory

This directory contains

- A docker file named like "Dockerfile.service"
- images of DTA services that can be deployed as containers

```shell
$> file executable
executable: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, Go BuildID=Jny1fdGJMrDqfG0rsISA/BwFdqPzR8hH3nq0RoRyF/Ks1i3mCYNWfAniZRZySr/t9I34Llb146kw2gh0_FD, not stripped
```

The relative small size of the docker-images is due to the fact, that only the statically linked executable is build into an empty docker-image (FROM scratch)

## Example dockerfile for `qds_count`

```Docker
FROM scratch
EXPOSE 50051
ADD "./docker/qds_count" /
CMD ["/qds_count"]
```

## Instantiation of a container

For executing the above image you could use

```shell
    docker run -p <exposedPort>:<hostPort> -d -it <container_name>
```
