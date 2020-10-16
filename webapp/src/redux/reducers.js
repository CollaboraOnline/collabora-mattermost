import {OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL} from './action_types';

export const rootModalData = (state, action) => {
    switch (action.type) {
    case OPEN_ROOT_MODAL:
        return {
            visible: true,
            fileId: action.payload.fileId,
        };
    case CLOSE_ROOT_MODAL:
        return {
            visible: false,
        };
    default:
        return typeof state === 'undefined' ? {visible: false} : state;
    }
};
