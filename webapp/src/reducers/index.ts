import {combineReducers} from 'redux';

import {wopiFilesList} from './wopi';
import {filePreviewModal} from './file_preview_modal';

export default combineReducers({
    wopiFilesList,
    filePreviewModal,
});
