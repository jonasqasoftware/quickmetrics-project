import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-metrics-list',
  templateUrl: './metrics-list.component.html',
  styleUrls: ['./metrics-list.component.css']
})
export class MetricsListComponent {
  @Input() metrics: any[];
}
