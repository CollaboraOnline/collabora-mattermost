import {AnyAction} from 'redux';

import Constants from '../constants';

export const config = (state = {}, action: AnyAction) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.RECEIVED_CLIENT_CONFIG:
        return action.data;
    case Constants.ACTION_TYPES.CLIENT_CONFIG_ERROR:
        return {};
    default:
        return state;
    }
};
