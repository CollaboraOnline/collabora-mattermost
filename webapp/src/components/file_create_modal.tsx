import React, {FC, useCallback, useEffect} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {Modal, Tooltip, OverlayTrigger, FormGroup, FormControl} from 'react-bootstrap';
import clsx from 'clsx';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import {closeFileCreateModal, createFileFromTemplate} from 'actions/file';
import {createFileModal} from 'selectors';

import {FILE_TEMPLATES, TEMPLATE_TYPES} from '../constants';

type FileCreateModalSelector = {
    visible: boolean;
    templateType: TEMPLATE_TYPES;
}

export const FileCreateModal: FC = () => {
    const {visible, templateType}: FileCreateModalSelector = useSelector(createFileModal);
    const currentChannelID: string = useSelector(getCurrentChannelId);

    const inputRef = React.createRef<HTMLInputElement>();
    const [fileName, setFileName] = React.useState('');
    const updateFileName = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFileName(e.target.value);
    };

    const [fileExt, setFileExt] = React.useState('');
    const updateFileExt = (e: React.ChangeEvent<FormControl>) => {
        setFileExt((e as unknown as React.ChangeEvent<HTMLSelectElement>).target.value);
    };

    // set the default extension when the template type changes
    useEffect(() => {
        setFileExt(FILE_TEMPLATES[templateType][0]);
    }, [templateType]);

    const onClear = useCallback((e?: React.MouseEvent<HTMLDivElement> | React.TouchEvent<HTMLDivElement>) => {
        e?.preventDefault();
        e?.stopPropagation();
        if (inputRef.current?.value) {
            inputRef.current.value = '';
        }
        setFileName('');
        inputRef.current?.focus();
    }, [inputRef]);

    const dispatch = useDispatch();
    const handleClose = useCallback((e?: React.MouseEvent<HTMLButtonElement>) => {
        e?.preventDefault?.();
        dispatch(closeFileCreateModal());
        onClear();
    }, [dispatch, onClear]);

    const handleConfirm = useCallback(() => {
        dispatch(createFileFromTemplate(currentChannelID, fileName, fileExt));
        handleClose();
    }, [fileExt, dispatch, currentChannelID, fileName, handleClose]);

    return (
        <Modal
            show={visible}
            onHide={handleClose}
            backdrop={'static'}
        >
            <Modal.Header closeButton={true}>
                <h4 className='modal-title'>
                    {`Create new ${templateType} from template.`}
                </h4>
            </Modal.Header>
            <Modal.Body>
                <div className='d-flex'>
                    <div className='form-control d-flex collabora-filename-container'>
                        <input
                            className='collabora-filename-input'
                            autoFocus={true}
                            type='text'
                            ref={inputRef}
                            maxLength={100}
                            value={fileName}
                            onChange={updateFileName}
                            placeholder={'Name your file'}
                        />
                        {
                            fileName && (
                                <div
                                    style={{right: 5}}
                                    className={clsx('input-clear', {visible: fileName.length > 0})}
                                    onMouseDown={onClear}
                                    onTouchEnd={onClear}
                                >
                                    <OverlayTrigger
                                        delayShow={400}
                                        placement={'bottom'}
                                        overlay={(
                                            <Tooltip id={'InputClearTooltip'}>
                                                {'Clear'}
                                            </Tooltip>
                                        )}
                                    >
                                        <span
                                            className='input-clear-x'
                                            aria-hidden='true'
                                        >
                                            <i className='icon icon-close-circle'/>
                                        </span>
                                    </OverlayTrigger>
                                </div>
                            )
                        }
                    </div>
                    <div>
                        <FormGroup controlId='formControlsSelect'>
                            <FormControl
                                className='collabora-file-ext-select'
                                componentClass='select'
                                placeholder='Select File Extension'
                                value={fileExt}
                                onChange={updateFileExt}
                            >
                                {
                                    FILE_TEMPLATES[templateType].map((item) => (
                                        <option
                                            key={item}
                                            value={item}
                                        >
                                            {`.${item}`}
                                        </option>
                                    ))
                                }
                            </FormControl>
                        </FormGroup>
                    </div>
                </div>
            </Modal.Body>
            <Modal.Footer>
                <button
                    type='button'
                    className='btn btn-link cancel'
                    onClick={handleClose}
                >
                    {'Cancel'}
                </button>
                <button
                    type='submit'
                    className={clsx('btn btn-primary confirm', {
                        disabled: !fileName,
                    })}
                    onClick={handleConfirm}
                    disabled={!fileName}
                >
                    {'Create'}
                </button>
            </Modal.Footer>
        </Modal>
    );
};

export default FileCreateModal;
