import { Component } from '@angular/core';
import { NavController, NavParams } from 'ionic-angular';
import { Http, Headers} from '@angular/http';

import { ChatroomPage } from './../chatroom/chatroom';
import { DataProvider } from '../../providers/data/data';
@Component({
  selector: 'page-home',
  templateUrl: 'home.html'
})
export class HomePage {
public checknol = 0;

  public dataChats = [
/*    {
      id: 1,
      room: 'Metode Penelitian dan Telaah Pustaka', 
      lastChatText: 'Minggu depan kuliah digabung dan masing-masing kelompok mempresentasikan proposal penelitian mereka', 
      lastChatTime: '15.00', 
      unread: '5'
    },
    {
      id: 2,
      room: 'Temu Kembali Informasi', 
      lastChatText: 'Kuliah penggantinya di Ruang Kuliah S2, jangan sampai telat ya', 
      lastChatTime: '13.00', 
      unread: ''
    },
    {
      id: 3,
      room: 'Manajemen Perangkat Lunak', 
      lastChatText: 'Presentasi proyek akan dilaksanakan tanggal 22 Desember 2017', 
      lastChatTime: '08.00', 
      unread: '99+'
    },
    {
      id: 4,
      room: 'Etika Komputasi', 
      lastChatText: 'Minggu ini ga ada kuliah ya', 
      lastChatTime: '06.00', 
      unread: '13'
    }*/
  ];

  constructor(public navCtrl: NavController, public navParams: NavParams, public data: DataProvider,  public http: Http) {
  }

  openChatroom($event, id) {
    this.navCtrl.push(ChatroomPage, id);
  }

  ionViewDidLoad(){
    this.getAllCGD();
  }

  getAllCGD(){
    let header= new Headers();
    header.append('Content-type', 'application/json' );  
    this.http.get(this.data.urlGetCGD, {withCredentials: true, headers:header})
      .subscribe(res => {
        let data = JSON.parse(res['_body']);
        if(data != null){
          for(let i = 0; i < data['datacgd'].length; i++){
            const CGD = {
                  id: data['datacgd'][i].idcgd,
                  room :  data['datamatakuliah'][i].kodematakuliah+ " "+ data['datamatakuliah'][i].namamatakuliah + "(" + data['datajadwalkuliah'][i].KelasParalel + ")",
                  unread : data['unread'][i],
                  semester : data['datamatakuliah'][i].semester
                }
            this.dataChats.push(CGD);
          }
        }

      })
  }

}
