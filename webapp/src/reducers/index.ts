import { combineReducers } from 'redux';

import { config } from './config';
import { createFileModal } from './create_file_modal';
import { filePreviewModal } from './file_preview_modal';
import { wopiFilesList } from './wopi';

export default combineReducers({
  config,
  wopiFilesList,
  filePreviewModal,
  createFileModal,
});
