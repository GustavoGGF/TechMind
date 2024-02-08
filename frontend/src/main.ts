import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/Pages/Login/app.config';
import { AppComponent } from './app/Pages/Login/app.component';

bootstrapApplication(AppComponent, appConfig).catch((err) =>
  console.error(err)
);
