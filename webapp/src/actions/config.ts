import type {AnyAction, Dispatch} from 'redux';
import type {ThunkAction, ThunkDispatch} from 'redux-thunk';
import type {ActionResult} from 'mattermost-redux/types/actions';
import type {GlobalState} from '@mattermost/types/store';

import Constants from '../constants';
import Client from 'client';

export const setConfig = (data: unknown) => {
    return async (dispatch: Dispatch) => {
        dispatch({
            type: Constants.ACTION_TYPES.RECEIVED_CLIENT_CONFIG,
            data,
        });
    };
};

export const getConfig = (): ThunkAction<Promise<ActionResult>, GlobalState, undefined, AnyAction> => {
    return async (dispatch: ThunkDispatch<GlobalState, undefined, AnyAction>) => {
        let data = null;
        try {
            data = await Client.getConfig();
            dispatch(setConfig(data));
        } catch (error) {
            dispatch({
                type: Constants.ACTION_TYPES.CLIENT_CONFIG_ERROR,
                error,
            });

            return {
                data,
                error,
            };
        }

        return {
            data,
            error: null,
        };
    };
};
