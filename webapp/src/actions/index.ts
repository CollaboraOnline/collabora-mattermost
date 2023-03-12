import { getConfig, setConfig } from './config';
import { closeFileCreateModal, createFileFromTemplate, showFileCreateModal } from './file';
import { closeFilePreview, showFilePreview } from './preview';
import { getCollaboraFileURL, getWopiFilesList } from './wopi';

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
