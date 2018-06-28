#!/usr/bin/env python3
# -*- coding: utf-8 -*-

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

import os
import subprocess
import sys
import shutil
import urllib.request
from pathlib import Path

SWAGGER_CODEGEN_CLI = "swagger-codegen-cli.jar"
SWAGGER_VER = "3.0.0-rc1"


def download_swagger(ver):
    """
    Download swagger

    @ver string the swagger version to download
    """
    url = ("http://central.maven.org"
           "/maven2/io/swagger/swagger-codegen-cli/{0}"
           "/swagger-codegen-cli-{0}.jar").format(ver)
    filename = SWAGGER_CODEGEN_CLI
    urllib.request.urlretrieve(url, filename)


if __name__ == "__main__":
    """
    Download the swagger codegen cli and run it on our api declaration file.
    """
    my_file = Path(SWAGGER_CODEGEN_CLI)
    if not my_file.is_file():
        print("Download swagger codegen cli")
        download_swagger(SWAGGER_VER)

    # Run the swagger codegen.
    print("Make a tmp directory for the generated server code")
    try:
        os.mkdir("./server")
    except FileExistsError as e:
        print("WARNING: \"server\" directory already exist.")

    # Test the swagger declaration file.
    print("Generate code")
    subprocess.call([
        'java',
        '-jar',
        SWAGGER_CODEGEN_CLI,
        'generate',
        '-l',
        'java',
        '-i',
        '../customers-service.yaml',
        '-o',
        './server'])

    # Remove the tmp directory.
    for arg in sys.argv:
        if arg == "-r":
            print("Delete tmp directory.")
            shutil.rmtree("./server")
