import {createAsyncAction} from 'typesafe-actions';

export const refreshToken = createAsyncAction(
  'TOKEN_REFRESH_REQUESTED',
  'TOKEN_REFRESH_SUCCEDED',
  'TOKEN_REFRESH_FAILED'
)<void, string, Error>();
