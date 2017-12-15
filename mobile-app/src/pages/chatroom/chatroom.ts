import { Component, ViewChild } from '@angular/core';
import { IonicPage, NavController, NavParams } from 'ionic-angular';
import { Events, Content, TextInput } from 'ionic-angular';

@IonicPage()
@Component({
  selector: 'page-chatroom',
  templateUrl: 'chatroom.html',
})
export class ChatroomPage {
  @ViewChild(Content) content: Content;
  @ViewChild('chat_input') messageInput: TextInput;
  chatInfo: any;
  userId = 1;
  msgList = [
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
  ];

  constructor(public navCtrl: NavController, public navParams: NavParams, public events: Events) {
    this.chatInfo = this.navParams;
  }

  onFocus() {
    this.content.resize();
    this.scrollToBottom();
  }

  sendMsg() {
    console.log('Message Send');
  }

  scrollToBottom() {
    setTimeout(() => {
      if (this.content.scrollToBottom) {
        this.content.scrollToBottom();
      }
    }, 400)
  }

}
