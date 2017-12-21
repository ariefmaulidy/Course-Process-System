import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router } from '@angular/router';
import { DataService } from '../data/data.service';
import swal from 'sweetalert2'


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
  private idPengajar: any[] = [];
  private startTime = {hour: 9, minute: 0};
  private endTime = {hour: 11, minute: 0};
  private pengajar: any;
  private topikkuliah: any;
  private jumlahpeserta: any;
  private catatan: any;
  private tanggal: any;
  private waktu: any;
  private jam_mulai: any;
  private menit_mulai: any;
  private jam_akhir: any;
  private menit_akhir: any;

  constructor(public http: Http, public router: Router, public dataService: DataService) { }

  ngOnInit() {
    window.scrollTo(0,0);
    this.matkul = "";
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
            this.dataJadwal[data['datamatakuliah'][i].kodematakuliah + " " + data['datamatakuliah'][i].namamatakuliah + " (" + data['datajadwalkuliah'][i].KelasParalel + ")"] = data['datajadwalkuliah'][i].idjadwalkuliah;
          }
        }

      })
  }

  getPengajar(namamatkul){
    this.http.get(this.dataService.urlGetPengajarBAP + '/' + this.dataMatkul[namamatkul],{withCredentials: true})
    .subscribe(res=>{
      let data = JSON.parse(res['_body']);
      this.namaPengajar = [];
      this.pengajar = null;
      if(data != null){
        for(let i  = 0; i < data.length; i++){
          this.namaPengajar[i] = data[i].nama;
          this.idPengajar[data[i].nama] = data[i].iduser;
        }
      }
    })
  }

  addBAP(){
    if(this.startTime['hour'] < 10){
      this.jam_mulai = '0' + this.startTime['hour'] 
    }else{
      this.jam_mulai = this.startTime['hour']
    }
    if(this.startTime['minute'] < 10){
      this.menit_mulai = '0' + this.startTime['minute'] 
    }else{
      this.menit_mulai = this.startTime['minute']
    }
    if(this.endTime['hour'] < 10){
      this.jam_akhir = '0' + this.endTime['hour'] 
    }else{
      this.jam_akhir = this.endTime['hour']
    }
    if(this.endTime['minute'] < 10){
      this.menit_akhir = '0' + this.endTime['minute'] 
    }else{
      this.menit_akhir = this.endTime['minute']
    }
    this.waktu = this.jam_mulai +':' + this.menit_mulai + '-' + this.jam_akhir + ':' + this.menit_akhir;
    let creds = JSON.stringify({tanggal: this.tanggal,
      idjadwalkuliah: this.dataJadwal[this.matkul],
      topikkuliah: this.topikkuliah,
      jumlahpeserta: Number(this.jumlahpeserta),
      iduser: this.idPengajar[this.pengajar],
      waktu: this.waktu,
      catatan: this.catatan
    });
    var headers = new Headers();
		headers.append("Content-Type", "application/json");
		this.http.post(this.dataService.urlAddBAP, creds, {withCredentials:true,headers: headers})
    .subscribe(res=>{
      if(res['status'] == 201){
				swal({
				  title : 'Berhasil ditambahkan',
				  type : 'success'
        });
        this.dataMatkul = [];
        this.dataJadwal = [];
        this.namaMatkul = [];
        this.idMatkul = [];
        this.idJadwalKuliah = [];
        this.namaPengajar = [];
        this.idPengajar = [];
        this.matkul = "";
        this.tanggal = "";
        this.topikkuliah = "";
        this.jumlahpeserta = "";
        this.waktu = "";
        this.catatan = "";
        this.pengajar = "";
				this.ngOnInit();
			} else if(res['_body'] == "duplicate"){
        console.log("gagal");
			} 
    });
  }

}
