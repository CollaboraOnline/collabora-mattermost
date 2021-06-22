import React, {FC} from 'react';
import {useSelector} from 'react-redux';
import clsx from 'clsx';
import {Button} from 'react-bootstrap';

import {FileInfo} from 'mattermost-redux/types/files';
import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';

import Client from 'client';

import CloseIcon from './close_icon';

type Props = {
    fileInfo: FileInfo;
    onClose: () => void;
    editable: boolean;
    toggleEditing: () => void;
}

export const FilePreviewHeader: FC<Props> = ({fileInfo, onClose, editable, toggleEditing}: Props) => {
    const post = useSelector((state: GlobalState) => getPost(state, fileInfo.post_id || ''));
    const channel = useSelector((state: GlobalState) => getChannel(state, post?.channel_id));
    let channelName: React.ReactNode = channel?.display_name || '';
    if (channel?.type === 'D') {
        channelName = 'Direct Message';
    } else if (channel?.type === 'G') {
        channelName = 'Group Message';
    }

    return (
        <>
            <div
                id='header'
                style={{
                    fontSize: 15,
                    lineHeight: 1.46668,
                    fontWeight: 400,
                    borderBottom: '1px solid #e1e1e1',
                    boxShadow: 'inset 0 1px 0 rgb(0 0 0 / 20%)',
                    height: 64,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    flex: '0 0 auto',
                }}
            >
                <div
                    id='headerMeta'
                    style={{
                        display: 'flex',
                        alignItems: 'center',
                        padding: '12px 16px',
                        minWidth: 0,
                    }}
                >
                    <div
                        style={{
                            maxHeight: 40,
                            minWidth: 0,
                        }}
                    >
                        <div
                            style={{
                                display: 'block',
                                textOverflow: 'ellipsis',
                                overflow: 'hidden',
                                whiteSpace: 'nowrap',
                                fontWeight: 700,
                            }}
                        >
                            {fileInfo.name}
                        </div>
                        <div
                            style={{
                                display: 'flex',
                                fontSize: 13,
                                lineHeight: 1.38463,
                                fontWeight: 400,
                            }}
                        >
                            <span
                                style={{
                                    color: '#606060',
                                    fontWeight: 700,
                                    paddingRight: 4,
                                    whiteSpace: 'nowrap',
                                    overflow: 'hidden',
                                    textOverflow: 'ellipsis',
                                }}
                            >
                                {channelName}
                            </span>
                        </div>
                    </div>
                </div>
                <div className='collabora-header-actions'>
                    <Button
                        bsSize='large'
                        bsStyle='large'
                        title='Download'
                        aria-label='Download'
                        className='collabora-header-action-button'
                        href={Client.getFileUrl(fileInfo.id)}
                        target='_blank'
                        rel='noopener noreferrer'
                        download={true}
                    >
                        <i className='fa fa-cloud-download'/>
                    </Button>
                    <Button
                        bsSize='large'
                        bsStyle='large'
                        onClick={toggleEditing}
                        className='collabora-header-action-button'
                        title={`${editable ? 'Lock' : 'Unlock'} Editing`}
                        aria-label={`${editable ? 'Lock' : 'Unlock'} Editing`}
                    >
                        <i
                            className={clsx(
                                'fa',
                                {
                                    'fa-lock': !editable,
                                    'fa-unlock': editable,
                                },
                            )}
                        />
                    </Button>
                    <div className='collabora-header-actions-separator'/>
                    <CloseIcon
                        id='closeIcon'
                        title='Close'
                        aria-label='Close'
                        className='close-x collabora-header-action-button'
                        onClick={onClose}
                    />
                </div>
            </div>
        </>
    );
};

export default FilePreviewHeader;
