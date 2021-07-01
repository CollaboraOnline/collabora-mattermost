import {combineReducers} from 'redux';

import {config} from './config';
import {wopiFilesList} from './wopi';
import {filePreviewModal} from './file_preview_modal';
import {createFileModal} from './create_file_modal';

export default combineReducers({
    config,
    wopiFilesList,
    filePreviewModal,
    createFileModal,
});
