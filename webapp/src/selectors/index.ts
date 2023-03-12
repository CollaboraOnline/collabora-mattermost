import { createSelector } from 'reselect';

import { FileInfo } from '@mattermost/types/files';
import { Post } from '@mattermost/types/posts';
import { GlobalState } from '@mattermost/types/store';
import { UserProfile } from '@mattermost/types/users';

import { getPost } from 'mattermost-redux/selectors/entities/posts';
import { getCurrentUser } from 'mattermost-redux/selectors/entities/users';

import { FILE_EDIT_PERMISSIONS } from '@/constants';
import { manifest } from '@/manifest';

// @ts-ignore GlobalState is not complete for plugins
const getPluginState = (state: GlobalState) => state[`plugins-${manifest.id}`] || {};

export const wopiFilesList = (state: GlobalState) => getPluginState(state).wopiFilesList;

export const filePreviewModal = (state: GlobalState) => getPluginState(state).filePreviewModal;

export const createFileModal = (state: GlobalState) => getPluginState(state).createFileModal;

export const collaboraConfig = (state: GlobalState) => getPluginState(state).config;

export const enableEditPermissions = (state: GlobalState) => Boolean(collaboraConfig(state)?.file_edit_permissions);

export function makeGetIsCurrentUserFileOwner(): (state: GlobalState, fileInfo: FileInfo) => boolean {
  return createSelector(
    'makeGetIsCurrentUserFileOwner',
    (state: GlobalState, fileInfo: FileInfo) => fileInfo,
    (state: GlobalState, fileInfo: FileInfo) => getPost(state, fileInfo.post_id || ''),
    (state: GlobalState) => getCurrentUser(state),
    (fileInfo: FileInfo, post: Post, currentUser: UserProfile) => {
      // for the existing attachment, user_id is fetched from post
      // but, for the newly created attachment, user_id is fetched from fileInfo,
      return Boolean(post?.user_id === currentUser.id || fileInfo.user_id === currentUser.id);
    }
  );
}

export function makeGetCollaboraFilePermissions(): (state: GlobalState, fileInfo: FileInfo) => FILE_EDIT_PERMISSIONS {
  return createSelector(
    'makeGetCollaboraFilePermissions',
    (state: GlobalState) => enableEditPermissions(state),
    (state: GlobalState, fileInfo: FileInfo) => getPost(state, fileInfo.post_id || ''),
    (state: GlobalState, fileInfo: FileInfo) => fileInfo.id,
    (featureEnabled: boolean, post: Post, fileID: string) => {
      if (!featureEnabled) {
        // if the feature id disabled, then everyone in the channel can edit
        return FILE_EDIT_PERMISSIONS.PERMISSION_CHANNEL;
      }

      return post?.props?.[`${manifest.id}_file_permissions_${fileID}`];
    }
  );
}
