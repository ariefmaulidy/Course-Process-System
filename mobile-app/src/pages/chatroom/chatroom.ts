import { Component, ViewChild , EventEmitter} from '@angular/core';
import { IonicPage, NavController, NavParams } from 'ionic-angular';
import { Events, Content, TextInput } from 'ionic-angular';
import { Http, Headers} from '@angular/http';

import { DataProvider } from '../../providers/data/data';

@IonicPage()
@Component({
  selector: 'page-chatroom',
  templateUrl: 'chatroom.html',
})
export class ChatroomPage {
  @ViewChild(Content) content: Content;
  @ViewChild('chat_input') messageInput: TextInput;
  chatInfo: any;
/*  userId = 1;*/
  /*msgList = [
    {
      messageId: 1,
      userId: 1,
      userName: 'Mahasiswa',
      time: '19.00',
      message: 'Kuliah pengganti besok jadi ga Pak?',     
      status: 'send'
    },
    {
      messageId: 2,
      userId: 2,
      userName: 'Dosen',
      time: '19.30',
      message: 'Jadi, ruangannya di RK X 3.01 ya',     
      status: 'send'
    },
    {
      messageId: 3,
      userId: 2,
      userName: 'Dosen',
      time: '19.31',
      message: 'Oiya jangan lupa tugas minggu lalu',     
      status: 'send'
    },
    {
      messageId: 4,
      userId: 1,
      userName: 'Mahasiswa',
      time: '19.40',
      message: 'Ok Pak',     
      status: 'pending'
    }
  ];*/

    //socket
  private ws: WebSocket;//server socket 
    private listener: EventEmitter<any> = new EventEmitter();
  //

  //isiDataSend//
  private namamatkul;
  private isipesan: any[] = [];
  private namapengirim: any[] = [];
  private classpengirim: any[] = [];
  private status: any[] = [];
  private passedtime: any[] = [];

/*  private dataChat: any[] = [];*/

  private checkStatus1 = "Bukan Saya";
  private checkStatus2 = "Saya"; 


  constructor(public navCtrl: NavController, public navParams: NavParams, public events: Events, public data: DataProvider,  public http: Http) {
    this.chatInfo = this.navParams;
     this.realTimeChatGroup(this.chatInfo)
  }

  onFocus() {
    this.content.resize();
    this.scrollToBottom();
  }

  sendMsg(pesan) {
    let creds = JSON.stringify({isipesan: pesan});

    var headers = new Headers();
    headers.append("Content-Type", "application/json");
    this.http.post(this.data.urlAddPesan + "/" +this.chatInfo, creds, {withCredentials: true,headers: headers})
    .subscribe(res => {
      if(res['status'] == 201){
        console.log("pesan dikirim");
      } else {
        console.log("pesan gagal dikirim");
      }
    });
  }

  scrollToBottom() {
    setTimeout(() => {
      if (this.content.scrollToBottom) {
        this.content.scrollToBottom();
      }
    }, 400)
  }


  //socket
  realTimeChatGroup(idcgd){
    this.ws = new WebSocket(this.data.urlSocket + "/" + idcgd);
    this.ws.onmessage = event => {
      var msg = JSON.parse(event.data);
      if(msg != null){
        this.listener.emit({"type": "message", "data":msg});
        for(var i=0; i < msg.length; i++){
            this.isipesan[i] = msg[i].isipesan;
            this.namapengirim[i] = msg[i].namapengirim;
            this.classpengirim[i] = msg[i].classpengirim;
            this.status[i] = msg[i].status;
            this.passedtime[i] = msg[i].passedtime;
        }
      }
    }
  }



}
