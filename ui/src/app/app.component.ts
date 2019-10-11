import { Component } from '@angular/core';
import { episodesApiService } from './services/episodesApi.service';
import { Episode } from './services/classes/episodes';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'tcaster';
  episodes: Episode[] = [];

  constructor(private _episodesApiService: episodesApiService) { }

  ngOnInit() {
    this._episodesApiService
      .getEpisodesList()
      .subscribe(data => { this.episodes = data })
  }
}
