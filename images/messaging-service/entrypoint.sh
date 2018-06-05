#!/bin/sh
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

# This script is the entry point for the ActiveMQ Artemis container.

# Run the server:
exec java \
  -Xms1G \
  -Xmx1G \
  -Xbootclasspath/a:$(ls /usr/share/artemis/lib/jboss-logmanager*.jar) \
  -classpath "/usr/share/artemis/lib/artemis-boot.jar" \
  -Dartemis.home="/usr/share/artemis" \
  -Dartemis.instance.etc="/etc/artemis" \
  -Dartemis.instance="/var/lib/artemis" \
  -Ddata.dir="/var/lib/artemis/data" \
  -Djava.io.tmpdir="/var/lib/artemis/tmp" \
  -Djava.library.path="/usr/share/artemis/bin/lib/linux-$(uname -m)" \
  -Djava.security.auth.login.config="/etc/artemis/login.config" \
  -Djava.util.logging.manager="org.jboss.logmanager.LogManager" \
  -Dlogging.configuration="file:/etc/artemis/logging.properties" \
  org.apache.activemq.artemis.boot.Artemis \
  run
