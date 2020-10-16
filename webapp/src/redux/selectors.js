import {id as pluginId} from '../manifest';

const getPluginState = (state) => state['plugins-' + pluginId] || {};

export const getRootModalData = (state) => {
    return getPluginState(state).rootModalData;
};
