import {delay, put, spawn, takeLatest} from 'redux-saga/effects';
import {ActionType, isActionOf} from 'typesafe-actions';

import {refreshToken} from './auth.actions';

const loginOrigin =
  process.env.NODE_ENV === 'development'
    ? process.env.SERVER_ORIGIN
    : document.location.origin;
const loginURL = `${loginOrigin}/api/auth/google/login`;
const refreshURL = `${document.location.origin}/api/auth/refresh`;

export default function* () {
  yield takeLatest(isActionOf(refreshToken.success), handleTokenChange);
  yield refreshOrRedirect();
  yield spawn(refreshRoutine);
}

function* refreshRoutine() {
  // TODO: take from auth service
  const TTL = 5 * 60 * 1000;

  while (true) {
    yield delay(TTL * 0.8);
    yield refreshOrRedirect();
  }
}

function* refreshOrRedirect() {
  try {
    yield put(refreshToken.request());
    const tokenResp: Response = yield fetch(refreshURL);

    if (tokenResp.status / 100 !== 2) {
      throw Error(tokenResp.statusText);
    }

    const token: string = yield tokenResp.json().then(x => x.auth_token);

    yield put(refreshToken.success(token));
  } catch (e) {
    if (e instanceof Error) {
      yield put(refreshToken.failure(e));
      window.location.replace(loginURL);
    }
  }
}

function handleTokenChange(action: ActionType<typeof refreshToken.success>) {
  // TODO: Update apollo client things here
  console.log(action.payload);

  localStorage.setItem('token', action.payload);
}
