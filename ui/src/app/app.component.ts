import { Component, ApplicationRef } from '@angular/core';
import { apiService } from './services/api.service';
import { Torrent } from './services/classes/torrents';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'undercast';
  torrentList: Torrent[] = [];
  torrent = new Torrent();
  isTorrentFormOpen = false;
  websocket: WebSocketSubject<any> = webSocket('ws://localhost:8080/api/ws');


  constructor(private _apiService: apiService) { }

  ngOnInit() {
    this._apiService.getTorrentsList().subscribe(data => {
      this.torrentList = data;
    });
    this.websocket.asObservable().subscribe(dataFromServer => {
      console.log(dataFromServer)
      for (const t of this.torrentList) {
        if (t.id == dataFromServer.id) {
          Object.assign(t, dataFromServer)
        }
      }
    });

  }

  async addTorrent() {
    const savedTorrent = await this._apiService.addTorrent(this.torrent);
    this.torrentList.push(savedTorrent);
    this.torrent = new Torrent();
    this.isTorrentFormOpen = false;
  }

  openTorrentForm() {
    this.isTorrentFormOpen = true
  }
}
