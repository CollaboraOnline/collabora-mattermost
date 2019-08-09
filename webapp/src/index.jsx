/* eslint-disable */
import { id as pluginId } from './manifest';
import PostType from './post_type';
import RootModal from './rootModal';
import { connect } from 'react-redux';
import { getRootModalData } from './redux/selectors';
import { combineReducers } from 'redux';
import { rootModalData } from './redux/reducers'

const mapStateToProps = (state) => ({
  modalData: getRootModalData(state),
});

export default class Plugin {
  initialize(registry, store) {
    registry.registerReducer(combineReducers({ rootModalData }));

    //modal with Collabora Online where the file will be viewed/edited
    registry.registerRootComponent(connect(mapStateToProps, null)(RootModal));

    //custom post type that appends files and actions (view/edit) at the end of the post with file
    registry.registerPostTypeComponent('custom_post_with_file', connect()(PostType));
  }
}

window.registerPlugin(pluginId, new Plugin());
