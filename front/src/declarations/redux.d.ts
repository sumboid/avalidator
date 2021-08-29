import {Action} from 'typesafe-actions';

export type RootAction = Action;

declare module 'typesafe-actions' {
  interface Types {
    RootAction: RootAction;
  }
}
