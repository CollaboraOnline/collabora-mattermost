import {AnyAction, Dispatch} from 'redux';
import {ThunkAction} from 'redux-thunk';

import {ActionResult} from 'mattermost-redux/types/actions';

import Client from '../client';
import Constants from '../constants';

export function getCollaboraFileURL(fileID: string) {
    return async () => {
        let data = null;
        try {
            data = await Client.getCollaboraOnlineURL(fileID);
        } catch (error) {
            return {data, error};
        }
        return {data, error: null};
    };
}

export function getWopiFilesList(): ThunkAction<Promise<ActionResult>, any, undefined, AnyAction> {
    return async (dispatch: Dispatch) => {
        let data = null;
        try {
            data = await Client.getWopiFilesList();
        } catch (error) {
            return {data, error};
        }
        dispatch({
            type: Constants.ACTION_TYPES.RECEIVED_WOPI_FILES_LIST,
            data,
        });
        return {data, error: null};
    };
}
