import React, {FC} from 'react';
import {useSelector} from 'react-redux';

import {Button} from 'react-bootstrap';

import {FileInfo} from 'mattermost-redux/types/files';
import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';

import Client from 'client';

import CloseIcon from './close_icon';

import './styles.css';

type Props = {
    fileInfo: FileInfo;
    onClose: () => void;
}

export const FilePreviewHeader: FC<Props> = ({fileInfo, onClose}: Props) => {
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
                <div
                    id='headerActions'
                    style={{
                        flex: '0 0 auto',
                        display: 'flex',
                        alignItems: 'center',
                        padding: 12,
                    }}
                >
                    <Button
                        bsSize='large'
                        bsStyle='large'
                        title='Download'
                        aria-label='Download'
                        className='collabora-action-button'
                        style={{
                            border: 0,
                            width: 40,
                            height: 40,
                            padding: 4,
                            borderRadius: 4,
                            marginBottom: 5,
                            alignItems: 'center',
                            display: 'inline-flex',
                            justifyContent: 'center',
                        }}
                        href={Client.getFileUrl(fileInfo.id)}
                        target='_blank'
                        rel='noopener noreferrer'
                        download={true}
                    >
                        <i className='fa fa-cloud-download'/>
                    </Button>
                    <div
                        id='separator'
                        style={{
                            width: 1,
                            height: 40,
                            backgroundColor: '#e1e1e1',
                            margin: '0 8px 6px',
                        }}
                    />
                    <CloseIcon
                        id='closeIcon'
                        title='Close'
                        aria-label='Close'
                        className='close-x'
                        onClick={onClose}
                        style={{
                            position: 'relative',
                            top: 'unset',
                            right: 'unset',
                            marginBottom: 5,
                        }}
                    />
                </div>
            </div>
        </>
    );
};

export default FilePreviewHeader;
