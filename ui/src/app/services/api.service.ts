import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Torrent } from './classes/torrents';

@Injectable()
export class apiService {
  constructor(private httpClient: HttpClient) {}

  getTorrentsList(): Observable<any> {
    return this.httpClient.get('/api/torrents ');
  }

  async addTorrent(torrentData: Torrent): Promise<Torrent> {
    const torrent = await this.httpClient
      .post('/api/torrents', torrentData)
      .toPromise();
    return torrent as Torrent;
  }
}
