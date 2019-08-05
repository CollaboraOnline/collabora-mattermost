/* eslint-disable */

import React from 'react';
import { CLOSE_ROOT_MODAL } from './action_types';
import { id as pluginId } from './manifest';

export default class RootModal extends React.Component {

    constructor(props) {
        super(props);
    }

    componentDidUpdate() {
        if (this.props.modalData.fileId == undefined)
            return;

        //ask the server for the Collabora Online URL where the file will be edited
        //and load it into the iframe
        const requestAddress = "/plugins/" + pluginId + "/collaboraURL?file_id=" + this.props.modalData.fileId;
        fetch(requestAddress).then((data) => data.json()).then((data) => {
            document.getElementById("collabora-submit-form").action=data.url;
            document.getElementById("collabora-form-access-token").value=data.access_token;
            document.getElementById("collabora-submit-form").submit();
        });
    }

    //document.getElementById("collabora-iframe").contentWindow.postMessage('{"MessageId":"Action_Save","SendTime":123,"Values":{"DontTerminateEdit":true,"DontSaveIfUnmodified":false}}',"http://localhost:9980")

    render() {
        if (!this.props.modalData.visible) {
            return null;
        }

        return (
            <div style={{ position: "fixed", top: 0, left: 0, right: 0, bottom: 0, backgroundColor: "rgba(0,0,0,0.8)", zIndex: 999, display: "flex", justifyContent: "center", alignItems: "center", overflow: "hidden" }}>
                <i
                    className='icon fa fa-times'
                    style={{ fontSize: '30px', position: 'absolute', top: '5px', right: "5px", color: "white", cursor: "pointer" }}
                    onClick={() => this.props.dispatch({ type: CLOSE_ROOT_MODAL })}
                />
                <div style={{ borderRadius: "10px", backgroundColor: "white", display: "inline-block" }}>
                    <form action="" method="POST" target="collabora-iframe" id='collabora-submit-form' >
                      <input id="collabora-form-access-token" name="access_token" value="" type="hidden"/>
                    </form>
                    <iframe height="640" width="800" id="collabora-iframe" name="collabora-iframe"></iframe>
                    <div style={{ textAlign: "center", padding: "3px" }}>
                        Powered by Collabora Online&nbsp;
                        <img style={{ height: "25px" }} src="https://www.collaboraoffice.com/wp-content/uploads/2019/03/collabora-productivity-nav-icon.png"></img>
                    </div>
                </div>
            </div>
        )
    }
}