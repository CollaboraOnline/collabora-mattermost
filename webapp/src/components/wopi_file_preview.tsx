import React, {FC, useEffect} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import Client from '../client';

interface ComponentProps {
    fileInfo: FileInfo;
}

export const WopiFilePreview: FC<ComponentProps> = (props: ComponentProps) => {
    useEffect(() => {
        const {fileInfo} = props;
        if (fileInfo?.id) {
            handleWopiFile(fileInfo.id);
        }
    }, [props.fileInfo?.id]);

    const handleWopiFile = async (fileID: string) => {
        //ask the server for the Collabora Online URL & token where the file will be edited
        //and load it into the iframe
        // TODO: Handle this API call failure
        const fileData = await Client.getCollaboraOnlineURL(fileID);

        //as the request to Collabora Online should be of POST type, a form is used to submit it.
        (document.getElementById('collabora-submit-form') as HTMLFormElement).action = fileData.url;
        (document.getElementById('collabora-form-access-token')as HTMLInputElement).value = fileData.access_token;
        (document.getElementById('collabora-submit-form') as HTMLFormElement).submit();
    };

    return (
        <div
            style={{
                overflowX: 'auto',
                overflowY: 'hidden',
                position: 'relative',
            }}
        >
            <form
                action=''
                method='POST'
                target='collabora-iframe'
                id='collabora-submit-form'
            >
                <input
                    id='collabora-form-access-token'
                    name='access_token'
                    value=''
                    type='hidden'
                />
            </form>
            <iframe
                height={window.innerHeight - 50}
                width={window.innerWidth - 200}
                id='collabora-iframe'
                name='collabora-iframe'
            />
        </div>
    );
};

export default WopiFilePreview;
