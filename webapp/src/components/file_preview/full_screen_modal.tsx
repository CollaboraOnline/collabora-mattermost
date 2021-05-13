import React from 'react';
import {CSSTransition} from 'react-transition-group';
import clsx from 'clsx';

import BackIcon from 'components/file_preview/back_icon';
import CloseIcon from 'components/file_preview/close_icon';

// This must be on sync with the animation time in mattermost-webapp's full_screen_modal.scss
const ANIMATION_DURATION = 100;

type Props = {
    compact?: boolean;
    show: boolean;
    onClose?: () => void;
    onGoBack?: () => void;
    children: React.ReactNode;
    ariaLabel?: string;
    ariaLabelledBy?: string;
};

class FullScreenModal extends React.PureComponent<Props> {
    private modal = React.createRef<HTMLDivElement>();

    componentDidMount() {
        document.addEventListener('keydown', this.handleKeypress);
        document.addEventListener('focus', this.enforceFocus, true);
        this.enforceFocus();
    }

    componentWillUnmount() {
        document.removeEventListener('keydown', this.handleKeypress);
        document.removeEventListener('focus', this.enforceFocus, true);
    }

    enforceFocus = () => {
        setTimeout(() => {
            const currentActiveElement = document.activeElement;
            if (this.modal && this.modal.current && !this.modal.current.contains(currentActiveElement)) {
                this.modal.current.focus();
            }
        });
    }

    handleKeypress = (e: KeyboardEvent) => {
        if (e.key === 'Escape' && this.props.show) {
            this.close();
        }
    }

    close = () => {
        return this.props.onClose?.();
    }

    render() {
        return (
            <CSSTransition
                in={this.props.show}
                classNames='FullScreenModal'
                mountOnEnter={true}
                unmountOnExit={true}
                timeout={ANIMATION_DURATION}
                appear={true}
            >
                <>
                    <div
                        className={clsx('FullScreenModal', {'FullScreenModal--compact': this.props.compact})}
                        ref={this.modal}
                        tabIndex={-1}
                        aria-modal={true}
                        aria-label={this.props.ariaLabel}
                        aria-labelledby={this.props.ariaLabelledBy}
                        role='dialog'
                    >
                        {this.props.onGoBack && (
                            <BackIcon
                                id='backIcon'
                                onClick={this.props.onGoBack}
                                className='back'
                                aria-label={'Back'}
                            />
                        )}
                        {this.props.onClose && (
                            <CloseIcon
                                id='closeIcon'
                                onClick={this.close}
                                className='close-x'
                                aria-label={'Close'}
                            />
                        )}
                        {this.props.children}
                    </div>
                    <div
                        tabIndex={0}
                        style={{display: 'none'}}
                    />
                </>
            </CSSTransition>
        );
    }
}

export default FullScreenModal;
