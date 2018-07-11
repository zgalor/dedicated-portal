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

import React, { Component } from 'react';
import classNames from 'classnames';
import { ListView, Button, Row, Col } from 'patternfly-react'
import PropTypes from 'prop-types'
import { CSSTransition, TransitionGroup} from 'react-transition-group'; 

export const renderActions = () => (
  <div>
    <Button>Details</Button>
  </div>
);

export const renderAdditionalInfoItems = itemProperties =>
  itemProperties &&
  Object.keys(itemProperties).map(prop => {
    const cssClassNames = classNames('pficon', {
      'pficon-flavor': prop === 'hosts',
      'pficon-cluster': prop === 'clusters',
      'pficon-container-node': prop === 'nodes',
      'pficon-image': prop === 'images'
    });
    return (
      <ListView.InfoItem key={prop}>
        <span className={cssClassNames} />
        <strong>{itemProperties[prop]}</strong> {prop}
      </ListView.InfoItem>
    );
});


class ClusterList extends Component {
  static propTypes = {
    clusters: PropTypes.array.isRequired,
  }

  shouldComponentUpdate(nextProps) {
    return (nextProps.clusters.length !== 0); 
  }
  
  render() {
    return (
      <div>
        <ListView>
        <TransitionGroup>
          {this.props.clusters.map(({ actions, properties, title, description, expandedContentText, hideCloseIcon }, index) => (
            <CSSTransition
              key={title}
              timeout={500}
              classNames="list"
              unmountOnExit>
              <ListView.Item
                actions={renderActions(actions)}
                checkboxInput={<input type="checkbox" />}
                leftContent={<ListView.Icon name="cluster" type="pf" />}
                additionalInfo={renderAdditionalInfoItems(properties)}
                heading={title}
                description={description}
                stacked={false}
                hideCloseIcon={false}>
                <Row>
                  <Col sm={11}>{expandedContentText}</Col>
                </Row>
              </ListView.Item>
            </CSSTransition>
          ))}
        </TransitionGroup>
          
        </ListView>

      </div>
    );
  }
}

export { ClusterList };
