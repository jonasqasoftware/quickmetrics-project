import { Component, OnInit } from '@angular/core';
import { MetricsService } from './services/metrics.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  metrics: any[];
  loading: boolean = true;

  constructor(private metricsService: MetricsService) {}

  ngOnInit() {
    this.metricsService.getMetrics()
      .subscribe((data: any) => {
        this.metrics = data;
        this.loading = false;
      });
  }
}
