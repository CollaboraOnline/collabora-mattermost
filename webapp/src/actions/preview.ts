import {Dispatch} from 'redux';

import {FileInfo} from 'mattermost-redux/types/files';

import Constants from '../constants';

export const showFilePreview = (fileInfo: FileInfo) => (dispatch: Dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.SHOW_FILE_PREVIEW,
        fileInfo,
    });
};

export const closeFilePreview = () => (dispatch: Dispatch) => {
    dispatch({
        type: Constants.ACTION_TYPES.CLOSE_FILE_PREVIEW,
    });
};
