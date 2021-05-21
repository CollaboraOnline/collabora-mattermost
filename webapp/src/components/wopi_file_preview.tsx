import React, {FC, useEffect, useState} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import Client from 'client';

type Props = {
    fileInfo: FileInfo;
}

export const WopiFilePreview: FC<Props> = (props: Props) => {
    const [windowWidth, setWindowWidth] = useState(window.innerWidth);
    const [windowHeight, setWindowHeight] = useState(window.innerHeight);

    const handleResize = () => {
        setWindowHeight(window.innerHeight);
        setWindowWidth(window.innerWidth);
    };

    useEffect(() => {
        window.addEventListener('resize', handleResize);
        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, []);

    const fileID = props.fileInfo?.id;
    useEffect(() => {
        if (fileID) {
            handleWopiFile(fileID);
        }
    }, [fileID]);

    const handleWopiFile = async (selectedFileID: string) => {
        //ask the server for the Collabora Online URL & token where the file will be edited
        //and load it into the iframe
        // TODO: Handle this API call failure
        const fileData = await Client.getCollaboraOnlineURL(selectedFileID);

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
                flex: '1 1 0',
                background: '#f6f6f6',
                borderTop: 'none',
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
                id='collabora-iframe'
                name='collabora-iframe'
                height={windowHeight - 69}
                width={windowWidth - 5}
                style={{border: 'none'}}
            />
        </div>
    );
};

export default WopiFilePreview;
