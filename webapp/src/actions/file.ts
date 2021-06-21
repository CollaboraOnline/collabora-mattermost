import {Dispatch} from 'redux';

import {DispatchFunc} from 'mattermost-redux/types/actions';

import Constants, {TEMPLATE_TYPES} from '../constants';
import Client from '../client';

export const showFileCreateModal = (templateType: TEMPLATE_TYPES) => (dispatch: Dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.SHOW_FILE_CREATE_MODAL,
        templateType,
    });
};

export const closeFileCreateModal = () => (dispatch: Dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.CLOSE_FILE_CREATE_MODAL,
    });
};

export function createFileFromTemplate(channelID: string, name: string, ext: string): DispatchFunc {
    return async () => {
        let data = null;
        try {
            data = await Client.createFileFromTemplate(channelID, name, ext);
        } catch (error) {
            return {data, error};
        }
        return {data, error: null};
    };
}
