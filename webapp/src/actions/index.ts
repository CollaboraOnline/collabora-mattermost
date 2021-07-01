import {getConfig, setConfig} from './config';
import {getWopiFilesList, getCollaboraFileURL} from './wopi';
import {showFilePreview, closeFilePreview} from './preview';
import {createFileFromTemplate, closeFileCreateModal, showFileCreateModal} from './file';

export default {
    getConfig,
    setConfig,
    showFilePreview,
    closeFilePreview,
    getCollaboraFileURL,
    getWopiFilesList,
    createFileFromTemplate,
    closeFileCreateModal,
    showFileCreateModal,
};
