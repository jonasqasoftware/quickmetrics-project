import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { MetricCardComponent } from './components/metric-card/metric-card.component';
import { MetricsListComponent } from './components/metrics-list/metrics-list.component';
import { MetricsService } from './services/metrics.service';

@NgModule({
  declarations: [
    AppComponent,
    MetricCardComponent,
    MetricsListComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule
  ],
  providers: [MetricsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
