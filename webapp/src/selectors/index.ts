import {GlobalState} from 'mattermost-redux/types/store';

import {id as pluginId} from '../manifest';

//@ts-ignore GlobalState is not complete
const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};

export const wopiFilesList = (state: GlobalState) => getPluginState(state).wopiFilesList;
