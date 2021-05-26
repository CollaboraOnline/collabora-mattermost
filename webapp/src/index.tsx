import React from 'react';
import {AnyAction, Store} from 'redux';
import {ThunkDispatch} from 'redux-thunk';

//@ts-ignore PluginRegistry doesn't have types yet
import {PluginRegistry} from 'mattermost-webapp/plugins/registry';

import {GlobalState} from 'mattermost-webapp/types/store';
import {FileInfo} from 'mattermost-redux/types/files';

import {showFilePreview} from 'actions/preview';
import {getWopiFilesList} from 'actions/wopi';
import {wopiFilesList} from 'selectors';
import Reducer from 'reducers';

import FilePreviewModal from 'components/file_preview_modal';
import WopiFilePreview from 'components/wopi_file_preview';

import {id as pluginId} from './manifest';

export default class Plugin {
    shouldShowPreview = (store: Store<GlobalState>, fileInfo: FileInfo) => {
        const state = store.getState();
        const wopiFiles = wopiFilesList(state);
        return Boolean(wopiFiles?.[fileInfo.extension]);
    }

    public initialize(registry: PluginRegistry, store: Store<GlobalState>): void {
        registry.registerReducer(Reducer);
        registry.registerRootComponent(FilePreviewModal);
        const dispatch: ThunkDispatch<GlobalState, undefined, AnyAction> = store.dispatch;
        dispatch(getWopiFilesList());
        registry.registerFilePreviewComponent(
            this.shouldShowPreview.bind(null, store),
            (props: {fileInfo: FileInfo}) => (
                <WopiFilePreview
                    fileInfo={props.fileInfo}
                    editable={false}
                />
            ),
        );

        // ignore if registerFileDropdownMenuAction method does not exist
        registry.registerFileDropdownMenuAction?.(
            this.shouldShowPreview.bind(null, store),
            'Edit with Collabora',
            (fileInfo: FileInfo) => dispatch(showFilePreview(fileInfo)),
        );
    }
}

// @ts-ignore
window.registerPlugin(pluginId, new Plugin());
