import React, {FC, useCallback, useEffect, useState} from 'react';

import {useDispatch} from 'react-redux';

import {FileInfo} from 'mattermost-redux/types/files';

import {getCollaboraFileURL} from 'actions/wopi';

type Props = {
    editable: boolean;
    fileInfo: FileInfo;
    setLoading?: (_: boolean) => void;
}

export const WopiFilePreview: FC<Props> = (props: Props) => {
    const dispatch = useDispatch();
    const [error, setError] = useState(false);
    const [loading, setLoadingState] = useState(false);

    const setLoading = useCallback((currentlyLoading) => {
        setLoadingState(currentlyLoading);
        props.setLoading?.(currentlyLoading);
    }, [props]);

    const handleWopiFile = useCallback(async (selectedFileID: string) => {
        //ask the server for the Collabora Online URL & token where the file will be edited
        //and load it into the iframe

        setLoading(true);
        setError(false);
        const dispatchResult = await dispatch(getCollaboraFileURL(selectedFileID) as any);
        if (dispatchResult.error) {
            setLoading(false);
            setError(true);
            return;
        }

        setLoading(false);
        setError(false);

        const fileData = dispatchResult.data as {url: string, access_token: string};

        //as the request to Collabora Online should be of POST type, a form is used to submit it.
        (document.getElementById('collabora-submit-form') as HTMLFormElement).action = fileData.url + (props.editable ? '/edit' : '');
        (document.getElementById('collabora-form-access-token') as HTMLInputElement).value = fileData.access_token;
        (document.getElementById('collabora-submit-form') as HTMLFormElement).submit();
    }, [dispatch, props.editable]);

    useEffect(() => {
        const fileID = props.fileInfo?.id;
        if (fileID) {
            handleWopiFile(fileID);
        }
    }, [handleWopiFile, props.fileInfo]);

    if (loading) {
        return (
            <span
                id='loadingSpinner'
                className={'wopi-loading-spinner'}
            >
                <i
                    className='fa fa-spinner fa-fw fa-pulse spinner'
                    title={'Loading Icon'}
                />
            </span>
        );
    }

    if (error) {
        return (
            <div className='alert wopi-error'>
                <i className='fa fa-warning wopi-error-icon'/>
                <div>{'We\'re sorry, a file preview is not available.'}</div>
                <div>{'Please download to view the file.'}</div>
            </div>
        );
    }

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
