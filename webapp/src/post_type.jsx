/* eslint-disable */

import React from 'react';
import { id as pluginId } from './manifest';
import { OPEN_ROOT_MODAL } from './redux/action_types';


const { formatText, messageHtmlToComponent } = window.PostUtils;
/**
 * PostType component modifies a post that contains files.
 * It appends file name and the action that can be performed on that file (view or edit)
 */
export default class PostType extends React.Component {

    constructor(props) {
        super(props);
        const post = { ...this.props.post };
        const message = post.message || '';
        const formattedText = messageHtmlToComponent(formatText(message));

        //ask the server to parse file IDs from this post
        //server returns a list of file names that can be edited with Collabora Online
        const requestAddress = "/plugins/" + pluginId + "/fileInfo";
        fetch(requestAddress, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(post.file_ids)
        }).then((data) => data.json()).then((data) => {
            data = data.filter((e) => e.id !== "");//a file that cannot be edited/viewed looks like this {id: "", name: "", extension: ""}
            this.setState({
                files: data
            });
        });

        this.state = {
            files: [],//contains ONLY the files that can be viewed/edited with Collabora Online
            formatedText: formattedText
        };
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

