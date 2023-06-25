import { DispatchFunc } from 'mattermost-redux/types/actions';

import Client from '@/client';
import Constants, { FILE_EDIT_PERMISSIONS, TEMPLATE_TYPES } from '@/constants';

import { AppDispatch } from '@/../types/store';

export const showFileCreateModal = (templateType: TEMPLATE_TYPES) => (dispatch: AppDispatch) => {
  dispatch({
    type: Constants.ACTION_TYPES.SHOW_FILE_CREATE_MODAL,
    templateType,
  });
};

export const closeFileCreateModal = () => (dispatch: AppDispatch) => {
  dispatch({
    type: Constants.ACTION_TYPES.CLOSE_FILE_CREATE_MODAL,
  });
};

export function createFileFromTemplate(channelID: string, name: string, ext: string): DispatchFunc {
  return async () => {
    let data = null;
    try {
      data = await Client.createFileFromTemplate(channelID, name, ext);
    } catch (error) {
      return { data, error };
    }
    return { data, error: null };
  };
}

export function updateFileEditPermission(fileID: string, permission: FILE_EDIT_PERMISSIONS): DispatchFunc {
  return async () => {
    let data = null;
    try {
      data = await Client.updateFileEditPermission(fileID, permission);
    } catch (error) {
      return { data, error };
    }
    return { data, error: null };
  };
}
