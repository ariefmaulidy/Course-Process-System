import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router, ActivatedRoute } from '@angular/router';
import { DataService } from '../data/data.service';
import { forEach } from '@angular/router/src/utils/collection';
import swal from 'sweetalert2'

@Component({
  selector: 'app-bap',
  templateUrl: './bap.component.html',
  styleUrls: ['./bap.component.css'],
  providers: [DataService]
})
export class BapComponent implements OnInit {

  public id:number;
  public sub: any;

  private databap: any[] = [];
  private datapeserta: any[] = [];
  private datapengajar: any[] = [];
  private datajadwal;
  private datamatkul;
  private namadosen: string;
  private ruangan: string;
  
  constructor(public http: Http, public route: ActivatedRoute, public router: Router, public dataService: DataService){}  

  ngOnInit() {
    window.scrollTo(0,0);
    this.sub = this.route.params.subscribe(params => {
       this.id = +params['idjadwal'];
    });
    this.getBAP()
  }

  getBAP(){
    this.http.get(this.dataService.urlGetJadwalBAP + '/' + this.id, {withCredentials: true})
    .subscribe(res=>{
      let data = JSON.parse(res['_body']);
      if(data != null){
        this.databap = data['databap'];
        this.datapengajar = data['datapengajar'];
        this.datapeserta = data['datapeserta'];
        this.datamatkul = data['datamatkul'];
        this.datajadwal = data['datajadwal'];
        this.namadosen = data['datadosen'];
        this.ruangan = data['ruangan'];
      }
    })
  }

}
