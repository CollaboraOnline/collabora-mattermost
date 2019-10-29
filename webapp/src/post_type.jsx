/* eslint-disable */
import PropTypes from 'prop-types';
import React from 'react';
import { id as pluginId } from './manifest';
import { OPEN_ROOT_MODAL } from './redux/action_types';

const { formatText, messageHtmlToComponent } = window.PostUtils;
/**
 * PostType component modifies a post that contains files.
 * It appends file name and the action that can be performed on that file (view or edit)
 */
export default class PostType extends React.Component {

    static propTypes = {
        post: PropTypes.object.isRequired,
        fileInfos: PropTypes.arrayOf(PropTypes.object),
    }

    static wopiFiles;

    constructor(props) {
        super(props);
        const post = { ...this.props.post };
        const fileInfos = { ...this.props.fileInfos };
        const message = post.message || '';
        const formattedText = messageHtmlToComponent(formatText(message));
        this.state = {
            files: [],//contains ONLY the files that can be viewed/edited with Collabora Online
            formatedText: formattedText
        };

        //prepare post files
        Object.keys(fileInfos).forEach((key) => {
            let file = fileInfos[key];
            if(PostType.wopiFiles[file.extension] != undefined){
                this.state.files.push({
                    name: file.name,
                    id: file.id,
                    action: PostType.wopiFiles[file.extension].Action
                })
            }
        })
    }

    //when the users clicks view or edit file
    fileAction(fileId) {
        this.props.dispatch({
            type: OPEN_ROOT_MODAL,
            payload: { fileId: fileId }
        });
    }

    render() {
        //prepare HTML with file name and file action (view/edit)
        var files = [];
        for (var i = 0; i < this.state.files.length; i++) {
            let file = this.state.files[i];
            files.push(
                <div key={i}>
                    <i>{file.name}</i>&nbsp;
                    <u style={{ cursor: "pointer" }} onClick={this.fileAction.bind(this,
                        file.id)}>
                        {file.action}
                    </u>
                </div>
            );
        }

        return (
            <div>
                {/* post text */}
                {this.state.formatedText}
                {/* post files */}
                <div>
                    {files}
                </div>
            </div>
        );
    }
}

