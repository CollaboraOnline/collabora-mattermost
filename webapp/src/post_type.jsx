/* eslint-disable */

import React from 'react';
import { id as pluginId } from './manifest';
import { OPEN_ROOT_MODAL } from './action_types';


const { formatText, messageHtmlToComponent } = window.PostUtils;

export default class PostType extends React.Component {

    constructor(props) {
        super(props);
        const post = { ...this.props.post };
        const message = post.message || '';
        const formattedText = messageHtmlToComponent(formatText(message));

        //ask the server to parse file IDs and give a list of file names that can be edited with Collabora Online
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
            files: [],//contains only the files that can be edited with Collabora Online
            formatedText: formattedText
        };
    }

    editDocument(fileId) {
        this.props.dispatch({
            type: OPEN_ROOT_MODAL,
            payload: { fileId: fileId }
        });
    }

    render() {
        var files = [];
        for (var i = 0; i < this.state.files.length; i++) {
            let file = this.state.files[i];
            files.push(
                <div key={i}>
                    <i>{file.name}</i>&nbsp;
                    <u style={{ cursor: "pointer" }} onClick={this.editDocument.bind(this,
                        file.id)}>
                        {file.action}
                    </u>
                </div>
            );
        }

        return (
            <div>
                {this.state.formatedText}
                <div>
                    {files}
                </div>
            </div>
        );
    }
}

