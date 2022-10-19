Run this script
This just unregisters the images. They don't get cleaned up until garbage collection runs
Manually run garbage collection: On the vm hosting the registry run docker exec registry /bin/registry garbage-collect /etc/docker/registry/config.yml