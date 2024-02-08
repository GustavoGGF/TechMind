import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/Pages/Login/app.component';
import { config } from './app/Pages/Login/app.config.server';

const bootstrap = () => bootstrapApplication(AppComponent, config);

export default bootstrap;
