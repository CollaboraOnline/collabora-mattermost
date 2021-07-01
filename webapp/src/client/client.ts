import qs from 'qs';

import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';
import {Options} from 'mattermost-redux/types/client4';

import {id as pluginId} from '../manifest';
import {FILE_EDIT_PERMISSIONS} from '../constants';

export default class Client {
    apiURL: string;
    baseURL: string;

    constructor() {
        this.apiURL = '/api/v4';
        this.baseURL = `/plugins/${pluginId}/api/v1`;
    }

    getConfig = () => {
        return this.doGet(`${this.baseURL}/config`);
    }

    updateFileEditPermission = (fileID: string, permission: FILE_EDIT_PERMISSIONS) => {
        const params = {permission};
        const url = `${this.baseURL}/files/${fileID}/access${this.buildQueryString(params)}`;
        return this.doPost(url);
    }

    createFileFromTemplate = (channelID: string, name: string, ext: string) => {
        const params = {name, ext};
        return this.doPost(`${this.baseURL}/channels/${channelID}/files/new${this.buildQueryString(params)}`);
    };

    getFileUrl = (fileID: string) => {
        return `${this.apiURL}/files/${fileID}`;
    };

    getWopiFilesList = () => {
        // fetch wopiFiles, a JSON with file extensions, actions (view/edit) and the Collabora Online URL where the action is done
        return this.doGet(this.baseURL + '/wopiFileList');
    }

    getCollaboraOnlineURL = (fileID: string) => {
        // fetch the Collabora Online URL & token where the file will be edited
        const params = {
            file_id: fileID,
        };
        const url = `${this.baseURL}/collaboraURL${this.buildQueryString(params)}`;
        return this.doGet(url);
    }

    doGet = async (url: string, headers: Record<string, string> = {}) => {
        const options = {
            method: 'get',
            headers,
        };
        return this.doFetch(url, options);
    }

    doPost = async (url: string, body?: BodyInit, headers: Record<string, string> = {}) => {
        const options = {
            method: 'post',
            body: JSON.stringify(body),
            headers,
        };
        return this.doFetch(url, options);
    }

    doDelete = async (url: string, body: BodyInit, headers: Record<string, string> = {}) => {
        const options = {
            method: 'delete',
            headers,
        };
        return this.doFetch(url, options);
    }

    doPut = async (url: string, body: BodyInit, headers: Record<string, string> = {}) => {
        const options = {
            method: 'put',
            body: JSON.stringify(body),
            headers,
        };
        return this.doFetch(url, options);
    }

    doFetch = async (url: string, options: Options = {}) => {
        const {data} = await this.doFetchWithResponse(url, options);

        return data;
    };

    doFetchWithResponse = async (url: string, options: Options = {}) => {
        const response = await fetch(url, Client4.getOptions(options));

        let data;
        if (response.ok) {
            data = await response.json();

            return {
                response,
                data,
            };
        }

        const text = await response.text();

        throw new ClientError(Client4.url, {
            message: text || '',
            status_code: response.status,
            url,
        });
    };

    buildQueryString(parameters: Record<string, string | number | boolean>) {
        if (Object.keys(parameters).length === 0) {
            return '';
        }

        return `?${qs.stringify(parameters, {
            encodeValuesOnly: true,
            arrayFormat: 'indices',
        })}`;
    }
}
