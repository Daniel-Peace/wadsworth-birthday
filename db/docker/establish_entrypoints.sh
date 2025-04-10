#!/bin/bash

# start SSH daemon
/usr/sbin/sshd

# start MongoDB
/usr/bin/mongod -f /etc/mongod.conf
