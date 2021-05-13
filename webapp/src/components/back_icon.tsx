import React, {ButtonHTMLAttributes} from 'react';

export default function BackIcon(props: ButtonHTMLAttributes<unknown>) {
    return (
        <button
            {...props}
        >
            <svg
                width='24px'
                height='24px'
                viewBox='0 0 24 24'
                role='icon'
                aria-label={'Back Icon'}
            >
                <path d='M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z'/>
            </svg>
        </button>
    );
}
