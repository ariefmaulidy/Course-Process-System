import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router } from '@angular/router';
import { DataService } from '../data/data.service';


@Component({
  selector: 'app-input-bap',
  templateUrl: './input-bap.component.html',
  styleUrls: ['./input-bap.component.css'],
  providers:[DataService]
})

export class InputBapComponent implements OnInit {
  private dataMatkul: any[] = [];
  private dataJadwal: any[] = [];
  private matkul : any;
  private idMatkul : any[] = [];
  private namaMatkul : any[] = [];
  private idJadwalKuliah : any[] = [];
  private namaPengajar: any[] = [];
  private startTime = {hour: 9, minute: 0};
  private endTime = {hour: 11, minute: 0};
  private pengajar: any;


  constructor(public http: Http, public router: Router, public dataService: DataService) { }

  ngOnInit() {
    window.scrollTo(0,0);
    this.matkul = " ";
    this.getDataMaktul();
  }

  getDataMaktul(){
    this.http.get(this.dataService.urlGetJadwalKuliah, {withCredentials: true})
      .subscribe(res => {
        let data = JSON.parse(res['_body']);
        if(data != null){
          for(let i = 0; i < data['datamatakuliah'].length; i++){
            // this.idMatkul.push(data['datamatakuliah'][i].idmatakuliah);
            this.namaMatkul.push(data['datamatakuliah'][i].kodematakuliah + " " + data['datamatakuliah'][i].namamatakuliah + " (" + data['datajadwalkuliah'][i].KelasParalel + ")");
            // this.idJadwalKuliah.push(data['datajadwalkuliah'].idjadwalkuliah);
            this.dataMatkul[data['datamatakuliah'][i].kodematakuliah + " " + data['datamatakuliah'][i].namamatakuliah + " (" + data['datajadwalkuliah'][i].KelasParalel + ")"] = data['datamatakuliah'][i].idmatakuliah;
            this.dataJadwal[data['datamatakuliah'][i].kodematakuliah + " " + data['datamatakuliah'][i].namamatakuliah + " (" + data['datajadwalkuliah'][i].KelasParalel + ")"] = data['datajadwalkuliah'].idjadwalkuliah;
          }
        }

      })
  }

  getPengajar(namamatkul){
    this.http.get(this.dataService.urlGetPengajarBAP + '/' + this.dataMatkul[namamatkul],{withCredentials: true})
    .subscribe(res=>{
      let data = JSON.parse(res['_body']);
      console.log(data)
      this.namaPengajar = [];
      this.pengajar = null;
      if(data != null){
        for(let i  = 0; i < data.length; i++){
          this.namaPengajar[i] = data[i].nama;
        }
      }
    })
  }

 /* addBAP(){
    let creds = JSON.stringify({matkul: })

  }*/

}
