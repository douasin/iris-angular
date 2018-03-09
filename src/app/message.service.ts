import { Injectable } from '@angular/core';

import { Observable } from 'rxjs/Observable';

import { Message } from './message';
import { BackendService } from './backend.service';

@Injectable()
export class MessageService {
    // File to handle apis of message service.

    private messageURL = 'api/message'; // URL to message api

    constructor(
        private backendService: BackendService
    ) { }

    postMessage(message: Message): Observable<any> {
        return this.backendService.getData(
            this.messageURL, message)
    }
}
