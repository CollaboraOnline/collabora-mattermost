import {Dictionary} from 'mattermost-redux/types/utilities';

import * as ACTION_TYPES from './action_types';

export enum TEMPLATE_TYPES {
    DOCUMENT = 'document',
    PRESENTATION = 'presentation',
    SPREADSHEET = 'spreadsheet',
}

export const FILE_TEMPLATES: Dictionary<string[]> = {
    [TEMPLATE_TYPES.DOCUMENT]: ['docx', 'odt'],
    [TEMPLATE_TYPES.PRESENTATION]: ['pptx', 'odp'],
    [TEMPLATE_TYPES.SPREADSHEET]: ['xlsx', 'ods'],
};

export const CHANNEL_TYPES = {
    CHANNEL_OPEN: 'O',
    CHANNEL_PRIVATE: 'P',
    CHANNEL_DIRECT: 'D',
    CHANNEL_GROUP: 'G',
};

export default Object.freeze({
    ACTION_TYPES,
    CHANNEL_TYPES,
    TEMPLATE_TYPES,
    FILE_TEMPLATES,
});
