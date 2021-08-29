import {StateType} from 'typesafe-actions';
import {combineReducers, createStore, applyMiddleware} from 'redux';
import createSagaMiddleware from 'redux-saga';

import authReducer from 'auth/auth.reducer';
import rootSaga from 'saga';

const rootReducer = combineReducers({
  auth: authReducer,
});

export default () => {
  const sagaMiddleware = createSagaMiddleware();
  const middleware = applyMiddleware(sagaMiddleware);
  const store = createStore(rootReducer, {}, middleware);

  sagaMiddleware.run(rootSaga);

  return store;
};

export type RootState = StateType<typeof rootReducer>;
