import { Component } from '@angular/core';
import { IonicPage, NavController, NavParams } from 'ionic-angular';

@IonicPage()
@Component({
  selector: 'page-jadwal',
  templateUrl: 'jadwal.html',
})
export class JadwalPage {
  public date = '2017-12-15';
  public dataMatkuls = [
    {nama: 'Manajemen Perangkat Lunak', ruang: 'Lab 2 Ilkom', mulai: '09.00', selesai: '11.00'},
    {nama: 'Temu Kembali Informasi', ruang: 'Ruang Kuliah S2 Ilkom', mulai: '13.30', selesai: '15.10'}
  ];

  constructor(public navCtrl: NavController, public navParams: NavParams) {
  }

}
