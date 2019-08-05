/* eslint-disable */
import { id as pluginId } from './manifest';
import PostType from './post_type';
import RootModal from './rootModal';
import { connect } from 'react-redux';
import { OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL } from './action_types';
import { getRootModalData } from 'selectors';
import { combineReducers } from 'redux';

const rootModalData = (state, action) => {
  switch (action.type) {
    case OPEN_ROOT_MODAL:
      return {
        visible: true,
        fileId: action.payload.fileId
      };
    case CLOSE_ROOT_MODAL:
      return {
        visible: false
      };
    default:
      return state===undefined?{visible:false}:state;
  }
};

const mapStateToProps = (state) => ({
  modalData: getRootModalData(state),
});

export default class Plugin {
  initialize(registry, store) {
    registry.registerReducer(combineReducers({ rootModalData }));
    registry.registerRootComponent(connect(mapStateToProps, null)(RootModal));
    registry.registerPostTypeComponent('custom_post_with_file', connect()(PostType));
  }
}

window.registerPlugin(pluginId, new Plugin());
