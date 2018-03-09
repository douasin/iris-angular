import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap } from 'rxjs/operators';

@Injectable()
export class BackendService {
    constructor(private http: HttpClient) { }

    private handleError<T> (operation = 'operation', result?: T) {
        return (error: any): Observable<T> => {
            console.error(error);
            console.log("something wrong");
            throw of(error);
            //return of(result as T); ignore error when getting data
        };
    }

    getData(url: string, data: object): Observable<any> {
        let tok = document.getElementById("csrf").getAttribute("value");
        // Set the csrf token to header which get from index.html
        let httpOptions = {
            headers: new HttpHeaders({
                'Content-Type': 'application/json',
                'X-Csrf-Token': tok,
            }),
            withCredentials: true
        }

        let bodyString = JSON.stringify(data);
        return this.http
        .post(url, bodyString, httpOptions)
        .pipe( 
            catchError(this.handleError<any>('getData'))
        );
    }
}
