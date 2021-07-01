import {useMemo} from 'react';
import {useSelector} from 'react-redux';

import {FileInfo} from 'mattermost-redux/types/files';
import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';

import {CHANNEL_TYPES} from '../constants';

// Hook
// Returns channelName for file
// Parameter used is the fileInfo
export const useChannelName = (fileInfo: FileInfo) => {
    const post = useSelector((state: GlobalState) => getPost(state, fileInfo.post_id || ''));
    const channel = useSelector((state: GlobalState) => getChannel(state, post?.channel_id));
    return useMemo(() => {
        if (!channel) {
            return '';
        }

        switch (channel.type) {
        case CHANNEL_TYPES.CHANNEL_DIRECT:
            return 'Direct Message';

        case CHANNEL_TYPES.CHANNEL_GROUP:
            return 'Group Message';

        default:
            return channel.display_name;
        }
    }, [channel]);
};
