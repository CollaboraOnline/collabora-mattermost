import React, {FC, useCallback, useState} from 'react';

import {Button} from 'react-bootstrap';

import {FileInfo} from 'mattermost-redux/types/files';

import WopiFilePreview from 'components/wopi_file_preview';

type Props = {
    fileInfo: FileInfo;
}

const FilePreviewComponent: FC<Props> = ({fileInfo}: Props) => {
    const [loading, setLoading] = useState(true);
    const [editable, setEditable] = useState(false);
    const enableEditing = useCallback(() => setEditable(true), []);
    return (
        <>
            <WopiFilePreview
                fileInfo={fileInfo}
                editable={editable}
                setLoading={setLoading}
            />
            {!loading && !editable && (
                <Button onClick={enableEditing}>
                    <span className='wopi-switch-to-edit-mode'>
                        <i className='fa fa-pencil-square-o'/>
                        {' Enable Editing'}
                    </span>
                </Button>
            )}
        </>
    );
};

export default FilePreviewComponent;
