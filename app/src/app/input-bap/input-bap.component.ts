import { Component, OnInit } from '@angular/core';


@Component({
  selector: 'app-input-bap',
  templateUrl: './input-bap.component.html',
  styleUrls: ['./input-bap.component.css']
})

export class InputBapComponent implements OnInit {
  private matkul : any;
  private daftarMatkul = ['pengbio', 'etikom', 'mppl'];
  private startTime = {hour: 9, minute: 0};
  private endTime = {hour: 11, minute: 0};

  constructor() { }

  ngOnInit() {
  }

}
