import React, {type ButtonHTMLAttributes} from 'react';
import {Button} from 'react-bootstrap';

export default function CloseIcon(props: ButtonHTMLAttributes<unknown>) {
    return (
        <Button
            type='button'
            {...props}
        >
            <svg
                width='24px'
                height='24px'
                viewBox='0 0 24 24'
                role='img'
                aria-label='Close Icon'
            >
                <path
                    fillRule='nonzero'
                    d='M18 7.209L16.791 6 12 10.791 7.209 6 6 7.209 10.791 12 6 16.791 7.209 18 12 13.209 16.791 18 18 16.791 13.209 12z'
                />
            </svg>
        </Button>
    );
}
