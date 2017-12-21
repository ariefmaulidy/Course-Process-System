import { Component, OnInit , EventEmitter} from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router , ActivatedRoute} from '@angular/router';
import { DataService } from '../data/data.service';

@Component({
  selector: 'app-chatgrup',
  templateUrl: './chatgrup.component.html',
  styleUrls: ['./chatgrup.component.css'],
  providers:[DataService]
})
export class ChatgrupComponent implements OnInit {

  constructor(public http: Http, public route: ActivatedRoute, public router: Router,public dataService: DataService) { }

  	//socket
	private ws: WebSocket;//server socket 
  	private listener: EventEmitter<any> = new EventEmitter();
	//

	//isiDataSend//
	private isipesan: any[] = [];
	private namapengirim: any[] = [];
	private classpengirim: any[] = [];
	private status: any[] = [];
	private passedtime: any[] = [];

/*  private dataChat: any[] = [];*/

  private checkStatus1 = "Bukan Saya";
  private checkStatus2 = "Saya"; 

    public id:number;
  public sub: any;

  private kirim;


  ngOnInit() {
    window.scrollTo(0,0);
      this.sub = this.route.params.subscribe(params => {
         this.id = +params['idcgd'];
      });
      this.realTimeChatGroup(this.id)
  }

  //socket
  realTimeChatGroup(idcgd){
  	this.ws = new WebSocket(this.dataService.urlSocket + "/" + idcgd);
  	this.ws.onmessage = event => {
  		var msg = JSON.parse(event.data);
      if(msg != null){
  			this.listener.emit({"type": "message", "data":msg});
  			for(var i=0; i < msg.length; i++){
    				this.isipesan[i] = msg[i].isipesan;
    				this.namapengirim[i] = msg[i].namapengirim;
            if (msg[i].status == this.checkStatus1){
              if (msg[i].classpengirim == "TataUsaha"){
                this.kirim = "http://placehold.it/50/55C1E7/fff&text=TU";
              } else if (msg[i].classpengirim == "Mahasiswa"){
                this.kirim = "http://placehold.it/50/43A047/fff&text=Mhs";
              } else if (msg[i].classpengirim == "Dosen"){
                this.kirim = "http://placehold.it/50/9C27B0/fff&text=Dosen";
              }
    				  this.classpengirim[i] = this.kirim;
    				} else {
              this.classpengirim[i] = msg[i].classpengirim;
            }
            this.status[i] = msg[i].status;
    				this.passedtime[i] = msg[i].passedtime;
  			}
  		}
  	}
  }

  sendPesan(pesan){
    let creds = JSON.stringify({isipesan: pesan});

    var headers = new Headers();
    headers.append("Content-Type", "application/json");
    this.http.post(this.dataService.urlAddPesan + "/" +this.id, creds, {withCredentials: true,headers: headers})
    .subscribe(res => {
      if(res['status'] == 201){
        console.log("pesan dikirim");
      } else {
        console.log("pesan gagal dikirim");
      }
    });
  }

  public getEventListener() {
    return this.listener;
	}

  ngOnDestroy() {
    this.sub.unsubscribe();
    this.ws.close();
  }

}
