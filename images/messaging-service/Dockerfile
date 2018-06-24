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

FROM centos:7

RUN \
    yum install -y \
        java-1.8.0-openjdk-headless \
        libaio \
        openssl \
        wget \
        && \
    yum clean all

RUN \
    wget --quiet --output-document=tarball "https://www.apache.org/dyn/closer.cgi?filename=activemq/activemq-artemis/2.6.1/apache-artemis-2.6.1-bin.tar.gz&action=download" && \
    echo b56d27107c6b362eb31a85d2a4720134b3142c5f2ed61d44a08eda57fc3764d6 tarball | sha256sum --check

RUN \
    mkdir --parents /usr/share/artemis && \
    tar --strip-components=1 --directory=/usr/share/artemis --extract --file=tarball && \
    rm tarball && \
    rm /usr/share/artemis/lib/artemis-amqp-protocol-*.jar && \
    rm /usr/share/artemis/lib/artemis-hornetq-protocol-*.jar && \
    rm /usr/share/artemis/lib/artemis-mqtt-protocol-*.jar && \
    rm /usr/share/artemis/lib/artemis-openwire-protocol-*.jar && \
    sed --in-place 's|href="console"|href="console/"|g' /usr/share/artemis/web/index.html

COPY \
    artemis-launcher \
    /usr/share/artemis/bin/launcher

EXPOSE \
    8161 \
    61613

CMD [ \
    "/usr/share/artemis/bin/launcher", \
    "--config=/etc/artemis", \
    "--instance=/var/lib/artemis" \
]

ENTRYPOINT [ \
    "/usr/share/artemis/bin/launcher" \
]
