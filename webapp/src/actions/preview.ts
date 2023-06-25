import { FileInfo } from '@mattermost/types/files';

import { AppDispatch } from '@/../types/store';
import Constants from '@/constants';

export const showFilePreview = (fileInfo: FileInfo) => (dispatch: AppDispatch) => {
  dispatch({
    type: Constants.ACTION_TYPES.SHOW_FILE_PREVIEW,
    fileInfo,
  });
};

export const closeFilePreview = () => (dispatch: AppDispatch) => {
  dispatch({
    type: Constants.ACTION_TYPES.CLOSE_FILE_PREVIEW,
  });
};
