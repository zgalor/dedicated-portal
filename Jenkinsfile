/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

pipeline {

  agent {
    node {
      // This label isn't really needed, but the 'label' is mandatory inside
      // 'node', so we need to use it in order to also set a custom workspace.
      // Note that this label has also to be assigned to all the Jenkins nodes
      // where this pipeline is inteded to run.
      label 'go'

      // The source needs to be in a directory that is in the Go source path,
      // otherwise the Go tools don't work correctly:
      customWorkspace "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/container-mgmt/dedicated-portal"
    }
  }

  environment {
    // Set the environment so that Go tools will work correctly:
    GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"

    // The build process uses 'which' to find the binaries installed by 'go
    // install', so we need to add the `GOBIN` directory:
    PATH = "${PATH}:${GOPATH}/bin"
  }

  stages {

    stage('build') {
      steps {
        sh '''
          # Build the images:
          make version=${GIT_COMMIT} tars
        '''
      }
    }

    stage('deploy') {
      steps {
        sshagent(['sandbox-admin-key']) {
          sh '''
            # Set the target user and host:
            SSH_USER="admin"
            SSH_HOST="sandbox-0.private"
            SSH_TARGET="${SSH_USER}@${SSH_HOST}"

            # Calculate the args for 'ssh' and 'scp':
            SSH_ARGS="-o StrictHostKeyChecking=no"

            # Copy the 'oc' command and configuration file:
            scp ${SSH_ARGS} ${SSH_TARGET}:/usr/bin/oc oc
            scp ${SSH_ARGS} ${SSH_TARGET}:.kube/config kubeconfig

            # Calculate the args for 'oc':
            OC_ARGS="--config=kubeconfig"

            # Remove unused images:
            ssh ${SSH_ARGS} ${SSH_TARGET} docker image prune --all --force
            ssh ${SSH_ARGS} ${SSH_TARGET} rm --force "*.tar"

            # Load the new images:
            rsync --rsh="ssh ${SSH_ARGS}" --compress *.tar ${SSH_TARGET}:.
            for tar in *.tar; do
              ssh ${SSH_ARGS} ${SSH_TARGET} docker load -i ${tar}
            done

            # Deploy the application:
            oc ${OC_ARGS} new-project dedicated-portal || true
            oc ${OC_ARGS} process \
              --filename=template.yml \
              --param=NAMESPACE=dedicated-portal \
              --param=VERSION=${GIT_COMMIT} \
              --param=DOMAIN=${SSH_HOST} \
              --param=PASSWORD=redhat123 \
            | \
            oc ${OC_ARGS} apply \
              --filename=-
          '''
        }
      }
    }

  }

  post {
    failure {
      slackSend color: "danger", message: "Deployment pipeline failed, see the details <${env.BUILD_URL}|here>."
    }

    always {
      cleanWs()
    }
  }

}
