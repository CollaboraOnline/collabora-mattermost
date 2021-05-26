import React, {FC, useCallback, useEffect} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import Client from 'client';

type Props = {
    editable: boolean;
    fileInfo: FileInfo;
}

export const WopiFilePreview: FC<Props> = (props: Props) => {
    const handleWopiFile = useCallback(async (selectedFileID: string) => {
        //ask the server for the Collabora Online URL & token where the file will be edited
        //and load it into the iframe
        // TODO: Handle this API call failure
        const fileData = await Client.getCollaboraOnlineURL(selectedFileID);

        //as the request to Collabora Online should be of POST type, a form is used to submit it.
        (document.getElementById('collabora-submit-form') as HTMLFormElement).action = fileData.url + (props.editable ? '/edit' : '');
        (document.getElementById('collabora-form-access-token')as HTMLInputElement).value = fileData.access_token;
        (document.getElementById('collabora-submit-form') as HTMLFormElement).submit();
    }, [props.editable]);

    useEffect(() => {
        const fileID = props.fileInfo?.id;
        if (fileID) {
            handleWopiFile(fileID);
        }
    }, [handleWopiFile, props.fileInfo]);

    return (
        <div className='wopi-iframe-container'>
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
                id='collabora-iframe'
                name='collabora-iframe'
                className='wopi-iframe'
            />
        </div>
    );
};

export default WopiFilePreview;
