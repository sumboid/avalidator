import React from 'react';
import ReactDOM from 'react-dom';
import {Provider, useSelector} from 'react-redux';
import {ApolloProvider} from '@apollo/client';

import createStore, {RootState} from 'store';
import apolloClient from 'apollo/client';

import RootPage from 'pages/root/root.page';

import '@fontsource/roboto';
import 'index.module.scss';

document.body.innerHTML += '<div id="root"/>';

const Root: React.FC = () => {
  const token = useSelector((state: RootState) => state.auth.token);

  return token ? (
    <ApolloProvider client={apolloClient}>
      <RootPage />
    </ApolloProvider>
  ) : (
    <div>Token loading...</div>
  );
};

ReactDOM.render(
  <React.StrictMode>
    <Provider store={createStore()}>
      <Root />
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);
