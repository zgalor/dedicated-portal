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
        wget \
        && \
    yum clean all

COPY \
    setup.sh \
    entrypoint.sh \
    etc \
    /tmp/

ADD \
    etc \
    /tmp/etc/

RUN \
    /tmp/setup.sh && \
    rm -rf /tmp/*

# 8161 - Administration console.
# 61613 - STOMP acceptor.
EXPOSE 8161 61613

CMD ["/root/entrypoint.sh"]
