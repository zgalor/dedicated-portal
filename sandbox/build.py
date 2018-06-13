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

import argparse
import jinja2
import os
import os.path
import subprocess
import sys

# The values extracted from the command line:
argv = None


def say(what):
    """
    Writes a message to the standard output, and then flushes it, so that
    the output doesn't appear out of order.
    """
    print(what, flush=True)


def generate_inventory():
    """
    Generates the Ansible inventory file from the template.
    """
    # Read and parse the template source:
    with open("inventory.j2", "r") as stream:
        source = stream.read()
    template = jinja2.Template(source)

    # Prepare a context containing all the command line options that have
    # a value:
    args = vars(argv)
    del args["func"]
    args = dict(
      (name, value) for name, value in args.items()
      if value is not None
    )

    # Render the template:
    result = template.render(args=args)
    with open("inventory.ini", "w") as stream:
        stream.write(result)


def generate_install_key():
    """
    Generates the temporary SSH keypair that will be used during
    the installation.
    """
    say("Generating temporary SSH key pair")
    for key_file in ["install_key", "install_key.pub"]:
        if os.path.exists(key_file):
            os.remove(key_file)
    process = subprocess.Popen([
        "ssh-keygen",
        "-b", "2048",
        "-t", "rsa",
        "-f", "./install_key",
        "-N", "",
    ])
    code = process.wait()
    if code != 0:
        raise Exception("SSH key generation failed with code {code}".format(
            code=code,
        ))


def run_playbook(path):
    """
    Runs the given Ansible playbook using the 'admin' SSH key pair.
    """
    say("Running playbook '{play}'".format(play=path))
    process = subprocess.Popen([
        "ansible-playbook",
        "--inventory", "inventory.ini",
        "--key-file", "install_key",
        path,
    ])
    code = process.wait()
    if code != 0:
        raise Exception("Playbook '{play}' failed with code {code}".format(
            play=path,
            code=code,
        ))


def add_sandbox_options(parser):
    """
    Adds the command line options that are used to configure the virtual
    machine.
    """
    parser.add_argument(
        "--sandbox-password",
        metavar="PASSWORD",
        help="The password of the 'admin' user.",
        default="redhat123",
    )
    parser.add_argument(
        "--sandbox-key",
        metavar="KEY",
        help=(
            "SSH public that will be added to the 'authorized_keys' file of "
            "the 'admin' user."
        ),
    )
    parser.add_argument(
        "--sandbox-hostname",
        metavar="HOSTNAME",
        help="The fully qualified host name of the sandbox virtual machine.",
        default="sandbox.local",
    )
    parser.add_argument(
        "--sandbox-address",
        metavar="IP",
        help="The IP address of the sandbox virtual machine.",
        default="192.168.122.100",
    )
    parser.add_argument(
        "--sandbox-netmask",
        metavar="MASK",
        help="The network mask of the sandbox virtual machine.",
        default="255.255.255.0",
    )
    parser.add_argument(
        "--sandbox-gateway",
        metavar="IP",
        help="The gateway the sandbox virtual machine.",
        default="192.168.122.1",
    )
    parser.add_argument(
        "--sandbox-dns-search",
        metavar="DOMAIN",
        help="The DNS search domain of the sandbox virtual machine.",
        default="local",
    )
    parser.add_argument(
        "--sandbox-dns-servers",
        metavar="SERVERS",
        help="The DNS servers of the sandbox virtual machine.",
        default="192.168.122.1",
    )


def add_libvirt_options(parser):
    """
    Adds the command line options that are specific to libvirt.
    """
    pass


def add_ovirt_options(parser):
    """
    Adds the command line options that are specifc to oVirt.
    """
    parser.add_argument(
        "--ovirt-url",
        metavar="URL",
        required=True,
        help="The URL of the oVirt API.",
    )
    parser.add_argument(
        "--ovirt-user",
        metavar="USER",
        required=True,
        help="The name of the oVirt user.",
        default="admin@internal",
    )
    parser.add_argument(
        "--ovirt-password",
        metavar="PASSWORD",
        required=True,
        help="The password of the oVirt user.",
    )
    parser.add_argument(
        "--ovirt-ca-file",
        metavar="FILE",
        required=True,
        help="The file containing trusted CA certificates.",
    )
    parser.add_argument(
        "--ovirt-template",
        metavar="TEMPLATE",
        required=True,
        help="The name of the oVirt template.",
    )
    parser.add_argument(
        "--ovirt-cluster",
        metavar="CLUSTER",
        required=True,
        help="The name of the oVirt cluster.",
    )


def build(env):
    """
    Builds the virtual machine using for the given kind of environment, which
    can be 'libvirt' or 'ovirt'.
    """
    generate_inventory()
    generate_install_key()
    run_playbook("clone_installer.yml")
    run_playbook("prepare_cloud_init.yml")
    run_playbook("prepare_{env}.yml".format(env=env))
    run_playbook("prepare_vm.yml")
    run_playbook("installer/playbooks/prerequisites.yml")
    run_playbook("installer/playbooks/deploy_cluster.yml")
    run_playbook("clean_vm.yml")


def build_libvirt():
    """
    Builds the virtual machine using libvirt.
    """
    build("libvirt")


def build_ovirt():
    """
    Builds the virtual machine using oVirt.
    """
    build("ovirt")


def main():
    # Create the top level command line parser:
    parser = argparse.ArgumentParser(
        prog=os.path.basename(sys.argv[0]),
        description="A simple tool to build the sandbox virtual machine.",
    )
    subparsers = parser.add_subparsers()

    # Create the parser for the 'libvirt' command:
    libvirt_parser = subparsers.add_parser("libvirt")
    add_sandbox_options(libvirt_parser)
    add_libvirt_options(libvirt_parser)
    libvirt_parser.set_defaults(func=build_libvirt)

    # Create the parser for the 'ovirt' command:
    ovirt_parser = subparsers.add_parser("ovirt")
    add_sandbox_options(ovirt_parser)
    add_ovirt_options(ovirt_parser)
    ovirt_parser.set_defaults(func=build_ovirt)

    # Run the selected tool:
    code = 0
    global argv
    argv = parser.parse_args()
    if hasattr(argv, "func"):
        argv.func()
    else:
        parser.print_usage()
        code = 1

    # Bye:
    sys.exit(code)


if __name__ == "__main__":
    main()
