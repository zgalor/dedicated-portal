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
import { connect } from 'react-redux';
import './App.css';
import * as fromItems from './ducks/items';
import { Header } from './Header.js';
import { ClusterList } from './ClusterList.js';
import { Pager, Label} from 'patternfly-react'
import PropTypes from 'prop-types'
import "patternfly/dist/css/patternfly.css";
import "patternfly/dist/css/patternfly-additions.css";

class App extends Component {
  componentDidMount() {
    const { fetchItems, itemsCurrentPage } = this.props;
    fetchItems(itemsCurrentPage);
    this.handleNext = this.handleNext.bind(this);
    this.handlePrevious = this.handlePrevious.bind(this);
  }
  handleNext() {
    const { fetchItems, itemsCurrentPage } = this.props;
    fetchItems(itemsCurrentPage + 1);
  }
  handlePrevious() {
    const { fetchItems, itemsCurrentPage } = this.props;
    fetchItems(itemsCurrentPage - 1);
  }
  render() {
    const {
      itemsPaged,
      itemsCurrentPage,
      itemsErrored,
      itemsLastPage,
      itemsRequested,
    } = this.props;
    let label;
    if (itemsRequested) label = <Label bsStyle="warning"> Requested </Label>;
    else if (itemsErrored) label = <Label bsStyle="danger"> Error fetching data </Label>
    else label = <Label bsStyle="success"> Updated </Label>


    let clusters = itemsPaged.map( item => Object.assign(item, {
      title: item.name,
      "properties": { "nodes": item['nodes'] },
      "expandedContentText":
        "Lorem Ipsum is simply dummy text of the printing and typesetting industry",
      "compoundExpandText": {
        "nodes": "Text describing Item 1s nodes"
      }}))
    return (
      <div>
        <Header></Header>
        {label}
        <ClusterList clusters={clusters}></ClusterList>
        <Pager
          messages={{nextPage: 'The Next Page', previousPage: 'The Previous Page'}}
          onNextPage={this.handleNext}
          onPreviousPage={this.handlePrevious}
        />
      </div>
    );
  }
}

App.propTypes = {
  fetchItems: PropTypes.func.isRequired,
  itemsPaged: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.number,
    name: PropTypes.string,
  })).isRequired,
  itemsCurrentPage: PropTypes.number.isRequired,
  itemsErrored: PropTypes.bool.isRequired,
  itemsLastPage: PropTypes.number.isRequired,
  itemsRequested: PropTypes.bool.isRequired,
};


const mapStateToProps = state => ({
  itemsCurrentPage: fromItems.getItemsCurrentPage(state),
  itemsErrored: fromItems.getItemsErrored(state),
  itemsLastPage: fromItems.getItemsLastPage(state),
  itemsPaged: fromItems.getItemsPaged(state),
  itemsRequested: fromItems.getItemsRequested(state),
});

const mapDispatchToProps = {
  fetchItems: fromItems.fetchItems,
};

export default connect(mapStateToProps, mapDispatchToProps)(App);