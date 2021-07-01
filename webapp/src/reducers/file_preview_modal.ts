import {AnyAction} from 'redux';

import Constants from '../constants';

const initialState = {
    visible: false,
    fileInfo: {},
    inhibited: false,
};

export const filePreviewModal = (state = initialState, action: AnyAction) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.SHOW_FILE_PREVIEW:
        return {
            ...state,
            visible: true,
            fileInfo: action.fileInfo,
        };

    // `inhibited` state allows other plugins to stop opening the full-screen collabora file preview component
    // and override that with a different component in that plugin.
    case Constants.ACTION_TYPES.INHIBIT_FILE_PREVIEW:
        return {
            ...state,
            inhibited: true,
        };

    case Constants.ACTION_TYPES.CLOSE_FILE_PREVIEW:
        return initialState;

    default:
        return state;
    }
};
