import { Component } from '@angular/core';
import { IonicPage, NavController, NavParams, AlertController } from 'ionic-angular';

@IonicPage()
@Component({
  selector: 'page-order',
  templateUrl: 'order.html',
})
export class OrderPage {
  private dataPesanans = [
    {
      ruangan: 'RK X 3.01',
      keperluan: 'Kuliah TKI',
      tanggal: 'Setiap Kamis',
      waktu: '08.00-09.40',
      pemesan: 'TU Ilkom',
      status: 'Lunas'
    },
    {
      ruangan: 'RK U 3.03',
      keperluan: 'Kuliah Pengganti Analgor',
      tanggal: '22 Desember 2017',
      waktu: '10.00-11.40',
      pemesan: 'TU Ilkom',
      status: 'Lunas'
    }
  ];

  constructor(public navCtrl: NavController, public navParams: NavParams, public alertCtrl: AlertController) {
  }

  showConfirm() {
    let confirm = this.alertCtrl.create({
      title: 'Setujui Pesanan?',
      message: 'Apakah Anda menyetujui pesanan ruangan ini?',
      buttons: [
        {
          text: 'Tidak',
          handler: () => {
            console.log('Disagree clicked');
          }
        },
        {
          text: 'Ya',
          handler: () => {
            console.log('Agree clicked');
          }
        }
      ]
    });
    confirm.present();
  }

}
