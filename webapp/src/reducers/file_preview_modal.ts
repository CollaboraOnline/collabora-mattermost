import {AnyAction} from 'redux';

import Constants from '../constants';

const initialState = {
    visible: false,
    fileInfo: {},
    editable: false,
};

export const filePreviewModal = (state = initialState, action: AnyAction) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.SHOW_FILE_PREVIEW:
        return {
            visible: true,
            fileInfo: action.fileInfo,
            editable: action.editable,
        };

    case Constants.ACTION_TYPES.CLOSE_FILE_PREVIEW:
        return initialState;

    default:
        return state;
    }
};
