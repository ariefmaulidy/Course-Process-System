import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Http,Headers } from '@angular/http';

/*
  Generated class for the DataProvider provider.

  See https://angular.io/guide/dependency-injection for more info on providers
  and Angular DI.
*/
@Injectable()
export class DataProvider {

	public hostname= "http://localhost:8080/"
	public hostWebsocket = 'ws://localhost:8080/';


  public urlGetJadwalMengajar = this.hostname + 'jadwalmengajar';
  public urlGetCGD = this.hostname + 'roomcgd'; 

  public urlAddPesan = this.hostname + 'addpesan'; 

  //socket
	public urlSocket = this.hostWebsocket + 'roomcgd';

  constructor(public http: HttpClient) {
    console.log('Hello DataProvider Provider');
  }


}
