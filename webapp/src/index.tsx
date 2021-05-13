import {Store} from 'redux';

//@ts-ignore Webapp imports don't work properly
import {PluginRegistry} from 'mattermost-webapp/plugins/registry';
import {GlobalState} from 'mattermost-webapp/types/store';

import {FileInfo} from 'mattermost-redux/types/files';

import {id as pluginId} from './manifest';

import FilePreviewOverride from './components/file_preview_override';
import FilePreviewModal from './components/file_preview/file_preview_modal';

import {getWopiFilesList} from './actions/wopi';
import {wopiFilesList} from './selectors';
import Reducer from './reducers';

export default class Plugin {
    public initialize(registry: PluginRegistry, store: Store<GlobalState>): void {
        registry.registerReducer(Reducer);
        registry.registerRootComponent(FilePreviewModal);
        registry.registerFilePreviewComponent(
            (fileInfo: FileInfo) => {
                const state = store.getState();
                const wopiFiles = wopiFilesList(state);
                return Boolean(wopiFiles?.[fileInfo.extension]);
            },
            FilePreviewOverride,
        );

        // @ts-ignore ThunkActions dont work properly
        store.dispatch(getWopiFilesList());
    }
}

// @ts-ignore
window.registerPlugin(pluginId, new Plugin());
