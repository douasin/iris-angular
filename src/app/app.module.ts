import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { NgModule } from '@angular/core';
import { HttpClientModule, HttpClient,
    HttpClientXsrfModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';

import { AppComponent } from './app.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';

import { MessageService } from './message.service'
import { BackendService } from './backend.service';
import { HomeComponent } from './home/home.component';


@NgModule({
  declarations: [
    AppComponent,
    PageNotFoundComponent,
    HomeComponent
  ],
  imports: [
      BrowserModule,
      FormsModule,
      AppRoutingModule,
      HttpClientModule,
      HttpClientXsrfModule.withOptions({
          // Those two settings need to be as same as iris did.
          cookieName: '_iris_csrf',
          headerName: 'X-Csrf-Token',
      }),
  ],
    providers: [
        MessageService,
        BackendService,
    ],
  bootstrap: [AppComponent]
})
export class AppModule { }
