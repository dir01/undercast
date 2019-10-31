import { Component } from '@angular/core';
import { apiService } from './services/api.service';
import { Torrent } from './services/classes/torrents';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'undercast';
  torrentList: Torrent[] = [];
  torrent = new Torrent();

  constructor(private _apiService: apiService) {}

  ngOnInit() {
    this._apiService.getTorrentsList().subscribe(data => {
      this.torrentList = data;
    });
  }

  async addTorrent() {
    const savedTorrent = await this._apiService.addTorrent(this.torrent);
    this.torrentList.push(savedTorrent);
    this.torrent = new Torrent();
  }
}
