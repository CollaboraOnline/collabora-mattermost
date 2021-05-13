import React, {FC, useCallback} from 'react';
import {useDispatch} from 'react-redux';

import {FileInfo} from 'mattermost-redux/types/files';

import {showFilePreview} from 'actions/preview';

type Props = {
    fileInfo: FileInfo;
}

export const FilePreviewOverride: FC<Props> = (props: Props) => {
    const dispatch = useDispatch();
    const openFilePreview = useCallback(() => {
        dispatch(showFilePreview(props.fileInfo));
    }, [dispatch, props.fileInfo]);

    return (
        <div>
            <button onClick={openFilePreview}>{'Collabora File Preview'}</button>
        </div>
    );
};

export default FilePreviewOverride;
