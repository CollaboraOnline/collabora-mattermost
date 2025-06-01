import React, {type FC, useCallback, useMemo, useState} from 'react';
import {useSelector} from 'react-redux';
import {Button} from 'react-bootstrap';

import type {FileInfo} from '@mattermost/types/files';
import type {GlobalState} from '@mattermost/types/store';

import WopiFilePreview from 'components/wopi_file_preview';
import {enableEditPermissions, makeGetCollaboraFilePermissions, makeGetIsCurrentUserFileOwner} from 'selectors';

import {FILE_EDIT_PERMISSIONS} from '../constants';

type Props = {
    fileInfo: FileInfo;
}

const FilePreviewComponent: FC<Props> = ({fileInfo}: Props) => {
    const [loading, setLoading] = useState(true);
    const [editable, setEditable] = useState(false);
    const enableEditing = useCallback(() => {
        setEditable(true);
    }, []);

    const getIsCurrentUserFileOwner = useMemo(makeGetIsCurrentUserFileOwner, []);
    const getCollaboraFilePermissions = useMemo(makeGetCollaboraFilePermissions, []);
    const isCurrentUserOwner = useSelector((state: GlobalState) => getIsCurrentUserFileOwner(state, fileInfo));
    const filePermission = useSelector((state: GlobalState) => getCollaboraFilePermissions(state, fileInfo));
    const editPermissionsFeatureEnabled = useSelector(enableEditPermissions);

    const showEditPermissionChangeOption = editPermissionsFeatureEnabled && isCurrentUserOwner;
    const canChannelEdit = filePermission === FILE_EDIT_PERMISSIONS.PERMISSION_CHANNEL;
    const canCurrentUserEdit = showEditPermissionChangeOption || canChannelEdit;

    return (
        <>
            <WopiFilePreview
                fileInfo={fileInfo}
                editable={editable}
                setLoading={setLoading}
            />
            {canCurrentUserEdit && !loading && !editable && (
                <Button onClick={enableEditing}>
                    <span className='wopi-switch-to-edit-mode'>
                        <i className='fa fa-pencil-square-o'/>
                        {' Edit'}
                    </span>
                </Button>
            )}
        </>
    );
};

export default FilePreviewComponent;
