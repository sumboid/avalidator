import React from 'react';
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';

import createStore from 'store';

import 'index.module.scss';

document.body.innerHTML += '<div id="root"/>';

ReactDOM.render(
  <React.StrictMode>
    <Provider store={createStore()}>
      Hello world!
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);
