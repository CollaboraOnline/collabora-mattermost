import React from 'react';
import {AnyAction, Store} from 'redux';
import {ThunkDispatch} from 'redux-thunk';

//@ts-ignore PluginRegistry doesn't have types yet
import {PluginRegistry} from 'mattermost-webapp/plugins/registry';

import {GlobalState} from 'mattermost-webapp/types/store';
import {FileInfo} from 'mattermost-redux/types/files';

import {showFileCreateModal} from 'actions/file';
import {showFilePreview} from 'actions/preview';
import {getWopiFilesList} from 'actions/wopi';
import {wopiFilesList} from 'selectors';
import Reducer from 'reducers';

import FilePreviewModal from 'components/file_preview_modal';
import FilePreviewComponent from 'components/file_preview_component';
import FileCreateModal from 'components/file_create_modal';

import {TEMPLATE_TYPES} from './constants';

import {id as pluginId} from './manifest';

import './components/styles.css';

export default class Plugin {
    shouldShowPreview = (store: Store<GlobalState>, fileInfo: FileInfo) => {
        const state = store.getState();
        const wopiFiles = wopiFilesList(state);
        return Boolean(wopiFiles?.[fileInfo.extension]);
    }

    public initialize(registry: PluginRegistry, store: Store<GlobalState>): void {
        registry.registerReducer(Reducer);
        registry.registerRootComponent(FilePreviewModal);
        registry.registerRootComponent(FileCreateModal);
        const dispatch: ThunkDispatch<GlobalState, undefined, AnyAction> = store.dispatch;
        dispatch(getWopiFilesList());
        registry.registerFilePreviewComponent(
            this.shouldShowPreview.bind(null, store),
            (props: {fileInfo: FileInfo}) => <FilePreviewComponent fileInfo={props.fileInfo}/>,
        );

        // ignore if registerFileDropdownMenuAction method does not exist
        registry.registerFileDropdownMenuAction?.(
            this.shouldShowPreview.bind(null, store),
            'Open with Collabora',
            (fileInfo: FileInfo) => dispatch(showFilePreview(fileInfo)),
        );

        registry.registerFileUploadMethod(
            <span className='fa wopi-file-upload-icon icon-filetype-document'/>,
            () => dispatch(showFileCreateModal(TEMPLATE_TYPES.DOCUMENT)),
            'New document',
        );
        registry.registerFileUploadMethod(
            <span className='fa wopi-file-upload-icon icon-filetype-spreadsheet'/>,
            () => dispatch(showFileCreateModal(TEMPLATE_TYPES.SPREADSHEET)),
            'New spreadsheet',
        );
        registry.registerFileUploadMethod(
            <span className='fa wopi-file-upload-icon icon-filetype-presentation'/>,
            () => dispatch(showFileCreateModal(TEMPLATE_TYPES.PRESENTATION)),
            'New presentation',
        );
    }
}

// @ts-ignore
window.registerPlugin(pluginId, new Plugin());
