import React, { ButtonHTMLAttributes } from 'react';
import { Button } from 'react-bootstrap';

export default function CloseIcon(props: ButtonHTMLAttributes<unknown>) {
  return (
    <Button type="button" {...props}>
      <svg width="24px" height="24px" viewBox="0 0 24 24" role="icon" aria-label="Close Icon">
        <path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z" />
      </svg>
    </Button>
  );
}
