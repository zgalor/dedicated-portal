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

# The name and version of the project:
project:=dedicated-portal
version:=latest

.PHONY: binaries
binaries: vendor
	for cmd in $$(ls cmd); do \
		cd cmd/$${cmd}; \
		go install || exit 1; \
		cd -; \
	done

vendor: Gopkg.lock
	dep ensure -vendor-only -v

.PHONY: images
images: binaries
	tmp=$$(mktemp -d); \
	trap "rm -rf $${tmp}" EXIT; \
	for image in $$(ls images); do \
		cp -r images/$${image}/* $${tmp}; \
		for cmd in $$(ls cmd); do \
			cp $$(which $${cmd}) $${tmp} || exit 1; \
		done; \
		tag=$(project)/$${image}:$(version); \
		docker build -t $${tag} $${tmp} || exit 1; \
	done

.PHONY: tars
tars: images
	for image in $$(ls images); do \
		tag=$(project)/$${image}:$(version); \
		tar=$$(echo $${tag} | tr /: __).tar; \
		docker save -o $${tar} $${tag} || exit 1; \
	done

.PHONY: tgzs
tgzs: tars
	for image in $$(ls images); do \
		tar=$(project)_$${image}_$(version).tar; \
		gzip -f $${tar} || exit 1; \
	done

.PHONY: lint
lint:
	./hack/verify.sh
	golint -min_confidence 0.9 -set_exit_status ./pkg/... ./cmd/...

.PHONY: fmt
fmt:
	gofmt -s -l -w ./pkg/ ./cmd/

.PHONY: clean
clean:
	rm -rf vendor *.tar *.tar.gz
