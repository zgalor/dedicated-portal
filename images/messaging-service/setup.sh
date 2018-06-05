#!/bin/sh -ex
#
# Copyright (c) 2018 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# This is the script that will run during the build of the container
# to install and configure the ActiveMQ Artemis server.

# Change to the directory where all the files have been copied:
cd /tmp

# Uncompress the tarball:
tar -xf artemis.tar.gz

# Remove the things that we don't want to install:
pushd apache-artemis-*
  rm -rf examples
  rm -rf web
  rm lib/artemis-amqp-protocol-*.jar
  rm lib/artemis-hornetq-protocol-*.jar
  rm lib/artemis-mqtt-protocol-*.jar
  rm lib/artemis-openwire-protocol-*.jar
popd

# Move the remaining files to their definitive location:
mkdir -p /usr/share/artemis
mv apache-artemis-*/* /usr/share/artemis/

# Create and populate the configuration directory:
mkdir -p /etc/artemis
mv etc/* /etc/artemis/.

# Move the start script to its location:
mv entrypoint.sh /root/
