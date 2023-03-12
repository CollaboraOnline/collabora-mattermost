import { ThunkDispatch } from 'redux-thunk';
import { GlobalState } from '@mattermost/types/store';
import { AnyAction } from 'redux';

export type AppDispatch = ThunkDispatch<GlobalState, undefined, AnyAction>;
