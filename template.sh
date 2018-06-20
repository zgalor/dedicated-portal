#!/bin/bash -ex

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

# This script uses the OpenShift template to create a deployment of
# the application.

# Create the namespace:
oc new-project dedicated-portal

# Use the template to create the objects:
oc process \
  --filename="template.yml" \
  --param=NAMESPACE="dedicated-portal" \
  --param=VERSION="latest" \
  --param=DOMAIN="example.com" \
  --param=PASSWORD="redhat123" \
| \
oc apply \
  --filename=-
