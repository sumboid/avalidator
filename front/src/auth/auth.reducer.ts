import {createReducer} from 'typesafe-actions';

import * as Actions from './auth.actions';

type State = {
  token: string | undefined;
  error: Error | undefined;
};

const DEFAULT_STATE: State = {
  token: undefined,
  error: undefined,
};

export default createReducer(DEFAULT_STATE)
  .handleAction(Actions.refreshToken.failure, (state, {payload: error}) => ({
    ...state,
    token: undefined,
    error,
  }))
  .handleAction(Actions.refreshToken.success, (state, {payload: token}) => ({
    ...state,
    token,
    error: undefined,
  }));
