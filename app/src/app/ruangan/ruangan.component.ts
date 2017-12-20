import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-ruangan',
  templateUrl: './ruangan.component.html',
  styleUrls: ['./ruangan.component.css']
})
export class RuanganComponent implements OnInit {
  private startTime = {hour: 9, minute: 0};
  private endTime = {hour: 11, minute: 0};
  
  constructor() { }

  ngOnInit() {
  }

}
