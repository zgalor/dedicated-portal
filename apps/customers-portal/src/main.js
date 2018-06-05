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

// Load styles:
import 'main.css'
import 'patternfly/dist/css/patternfly.min.css'
import 'patternfly/dist/css/patternfly-additions.min.css'

// Load dependencies:
import 'bootstrap/dist/js/bootstrap'
import 'patternfly/dist/js/patternfly'

import React from 'react'
import { render } from 'react-dom'

import App from 'App'

// Initialize Patternfly vertical navigation when the document is ready:
$(document).ready(() => {
  $().setupVerticalNavigation(true)
})

// Render the application:
render(
  <App/>,
  document.getElementById('app')
)
