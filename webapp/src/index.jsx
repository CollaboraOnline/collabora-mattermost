/* eslint-disable */
import { id as pluginId } from './manifest';
import PostType from './post_type';
import RootModal from './rootModal';
import { connect } from 'react-redux';
import { getRootModalData } from './redux/selectors';
import { combineReducers } from 'redux';
import { rootModalData } from './redux/reducers'
import {makeGetFilesForPost} from 'mattermost-redux/selectors/entities/files';

//for modal
const mapStateToProps = (state) => ({
  modalData: getRootModalData(state),
});

//for post file infos
function makeMapStateToProps() {
  const selectFilesForPost = makeGetFilesForPost();
  return function mapStateToProps(state, ownProps) {
      const postId = ownProps.post ? ownProps.post.id : '';
      const fileInfos = selectFilesForPost(state, postId);
      return {
          fileInfos,
      };
  };
}

export default class Plugin {
  initialize(registry, store) {
    registry.registerReducer(combineReducers({ rootModalData }));

    //modal with Collabora Online where the file will be viewed/edited
    registry.registerRootComponent(connect(mapStateToProps, null)(RootModal));

    //fetch wopiFiles, a JSON with file extensions, actions (view/edit) and the Collabora Online URL where the action is done
    fetch("/plugins/" + pluginId + "/wopiFileList")
      .then((data) => data.json())
      .then((data) => {
        PostType.wopiFiles = data;
        
        //custom post type that appends files and actions (view/edit) at the end of the post with file
        registry.registerPostTypeComponent('custom_post_with_file', connect(makeMapStateToProps)(PostType));
      });
  }
}

window.registerPlugin(pluginId, new Plugin());
