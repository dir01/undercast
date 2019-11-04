import { NgModule } from '@angular/core';
import {
  MatInputModule,
  MatProgressBarModule,
  MatCardModule,
  MatToolbarModule,
  MatDividerModule,
  MatIconModule
} from '@angular/material';
import { MatButtonModule } from '@angular/material/button';

const MaterialComponents = [
  MatInputModule,
  MatProgressBarModule,
  MatCardModule,
  MatToolbarModule,
  MatButtonModule,
  MatDividerModule,
  MatIconModule
]

@NgModule({
  imports: MaterialComponents,
  exports: MaterialComponents
})
export class MaterialModule { }
