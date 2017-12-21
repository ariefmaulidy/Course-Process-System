import { Component } from '@angular/core';
import { IonicPage, NavController, NavParams } from 'ionic-angular';

@IonicPage()
@Component({
  selector: 'page-ruangan-detail',
  templateUrl: 'ruangan-detail.html',
})
export class RuanganDetailPage {
  private item;
  private dataRuangans = [
    {
      keperluan: 'Kuliah TKI',
      tanggal: 'Setiap Kamis',
      waktu: '08.00-09.40',
      pemesan: 'TU Ilkom',
      status: 'Lunas'
    },
    {
      keperluan: 'Kuliah Pengganti Analgor',
      tanggal: '22 Desember 2017',
      waktu: '10.00-11.40',
      pemesan: 'TU Ilkom',
      status: 'Lunas'
    }
  ];

  constructor(public navCtrl: NavController, public navParams: NavParams) {
    this.item = this.navParams.data;
  }

}
