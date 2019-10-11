import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable()
export class episodesApiService {

    constructor(private httpClient: HttpClient) { }

    getEpisodesList(): Observable<any> {
        return this.httpClient.get('http://localhost:8080/episodes ')
    }
}