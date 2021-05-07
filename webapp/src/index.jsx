import {id as pluginId} from './manifest';

import WopiFilePreview from './components/wopi_file_preview';

import {getWopiFilesList} from './actions/wopi';
import {wopiFilesList} from './selectors';
import Reducer from './reducers';

export default class Plugin {
    initialize(registry, store) {
        registry.registerReducer(Reducer);
        registry.registerFilePreviewComponent(
            (fileInfo) => {
                const state = store.getState();
                const wopiFiles = wopiFilesList(state);
                return Boolean(wopiFiles?.[fileInfo.extension]);
            },
            WopiFilePreview,
        );
        store.dispatch(getWopiFilesList());
    }
}

window.registerPlugin(pluginId, new Plugin());
