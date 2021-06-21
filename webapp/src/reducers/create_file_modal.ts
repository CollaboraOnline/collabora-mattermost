import {AnyAction} from 'redux';

import Constants, {TEMPLATE_TYPES} from '../constants';

const initialState = {
    visible: false,
    templateType: TEMPLATE_TYPES.DOCUMENT,
};

export const createFileModal = (state = initialState, action: AnyAction) => {
    switch (action.type) {
    case Constants.ACTION_TYPES.SHOW_FILE_CREATE_MODAL:
        return {
            visible: true,
            templateType: action.templateType,
        };

    case Constants.ACTION_TYPES.CLOSE_FILE_CREATE_MODAL:
        return {
            ...state,
            visible: false,
        };

    default:
        return state;
    }
};
