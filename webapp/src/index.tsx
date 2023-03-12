import React from 'react';
import { Action, Store } from 'redux';

import { GlobalState } from '@mattermost/types/store';
import { FileInfo } from '@mattermost/types/files';
import { WebSocketMessage } from 'mattermost-redux/types/websocket';

import { PluginRegistry } from '@/../types/mattermost-webapp';
import { AppDispatch } from '@/../types/store';

import Actions from '@/actions';
import Reducer from '@/reducers';
import { wopiFilesList } from '@/selectors';

import FileCreateModal from '@/components/file_create_modal';
import FilePreviewComponent from '@/components/file_preview_component';
import FilePreviewModal from '@/components/file_preview_modal';

import Constants, { TEMPLATE_TYPES } from '@/constants';
import { manifest } from '@/manifest';

import './styles/styles.css';

export default class Plugin {
  shouldShowPreview = (store: Store<GlobalState>, fileInfo: FileInfo) => {
    const state = store.getState();
    const wopiFiles = wopiFilesList(state);
    return Boolean(wopiFiles?.[fileInfo.extension]);
  };

  public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
    // @see https://developers.mattermost.com/extend/plugins/webapp/reference/

    registry.registerReducer(Reducer);
    registry.registerRootComponent(FilePreviewModal);
    registry.registerRootComponent(FileCreateModal);

    const dispatch: AppDispatch = store.dispatch;
    dispatch(Actions.getConfig());
    registry.registerWebSocketEventHandler(
      `custom_${manifest.id}_${Constants.WEBSOCKET_EVENT_CONFIG_UPDATED}`,
      (event: WebSocketMessage<unknown>) => {
        dispatch(Actions.setConfig(event.data));
      }
    );

    dispatch(Actions.getWopiFilesList());
    registry.registerFilePreviewComponent(this.shouldShowPreview.bind(null, store), (props: { fileInfo: FileInfo }) => (
      <FilePreviewComponent fileInfo={props.fileInfo} />
    ));

    // ignore if registerFileDropdownMenuAction method does not exist
    registry.registerFileDropdownMenuAction?.(this.shouldShowPreview.bind(null, store), 'Open with Collabora', (fileInfo: FileInfo) =>
      dispatch(Actions.showFilePreview(fileInfo))
    );

    registry.registerFileUploadMethod(
      <span className="fa wopi-file-upload-icon icon-filetype-document" />,
      () => dispatch(Actions.showFileCreateModal(TEMPLATE_TYPES.DOCUMENT)),
      'New document'
    );
    registry.registerFileUploadMethod(
      <span className="fa wopi-file-upload-icon icon-filetype-spreadsheet" />,
      () => dispatch(Actions.showFileCreateModal(TEMPLATE_TYPES.SPREADSHEET)),
      'New spreadsheet'
    );
    registry.registerFileUploadMethod(
      <span className="fa wopi-file-upload-icon icon-filetype-presentation" />,
      () => dispatch(Actions.showFileCreateModal(TEMPLATE_TYPES.PRESENTATION)),
      'New presentation'
    );
  }
}

declare global {
  interface Window {
    registerPlugin(pluginId: string, plugin: Plugin): void;
  }
}

window.registerPlugin(manifest.id, new Plugin());
