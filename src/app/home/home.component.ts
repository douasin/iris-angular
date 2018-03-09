import { Component, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';

import { Message } from '../message';
import { MessageService } from '../message.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
    sending: boolean = false;
    success: boolean = false;
    sended: boolean = false;
    newMessage: Message = {
        name: "tester",
        email: "test@test.com",
        subject: "This is a test",
        content: "I'm content"
    };

    constructor(private messageService: MessageService) { }

    ngOnInit() { }

    reset(): void {
        this.sending = false;
    }

    sendMessage(): void {
        this.sended = false;
        this.sending = true;

        this.messageService.postMessage(this.newMessage)
        .subscribe(
            data => {
                console.log('send success');
                console.log(data);
                this.sended = true;
                this.success = true;
                this.reset();
            },
            (err: HttpErrorResponse) => {
                console.log('send fail');
                console.log(err);
                this.sended = true;
                this.sending = false;
                this.success = false;
            }
        );
    }
}
