// eslint-disable-next-line import/no-unresolved
import {ReadableStream} from 'node:stream/web';
// eslint-disable-next-line import/no-unresolved
import {Blob} from 'node:buffer';
// eslint-disable-next-line import/no-unresolved
import {MessagePort} from 'node:worker_threads';

class DOMException extends Error {
}

global.ReadableStream = ReadableStream;
global.Blob = Blob;
global.MessagePort = MessagePort;
global.DOMException = DOMException;
