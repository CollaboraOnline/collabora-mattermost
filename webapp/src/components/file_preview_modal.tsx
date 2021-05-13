import React, {FC, useCallback} from 'react';

import {FileInfo} from 'mattermost-redux/types/files';

import {useDispatch, useSelector} from 'react-redux';

import {closeFilePreview} from 'actions/preview';
import {filePreviewModal} from 'selectors';
import FullScreenModal from 'components/full_screen_modal';
import WopiFilePreview from 'components/wopi_file_preview';

const FilePreviewModal: FC = () => {
    const dispatch = useDispatch();
    const visible = useSelector(filePreviewModal)?.visible;
    const fileInfo: FileInfo = useSelector(filePreviewModal)?.fileInfo;

    const handleClose = useCallback((e?: Event): void => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }

        dispatch(closeFilePreview());
    }, [dispatch]);

    return (
        <FullScreenModal
            compact={true}
            show={visible}
            onClose={handleClose}
        >
            <WopiFilePreview fileInfo={fileInfo}/>
        </FullScreenModal>
    );
};

export default FilePreviewModal;
