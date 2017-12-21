import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router } from '@angular/router';
import { DataService } from '../data/data.service';
import swal from 'sweetalert2'



@Component({
  selector: 'app-ruangan',
  templateUrl: './ruangan.component.html',
  styleUrls: ['./ruangan.component.css'],
  providers: [DataService]
})
export class RuanganComponent implements OnInit {
  private dataMatkul: any[] = [];
  private dt;
  private tanggal;
  private dataruangan: any[] = [];
  private dataMataKuliah: any[] = [];
  private halaman: number = 1;
  private idruangan: any;
  private matkul : any;
  private alasan : any;
  private tglpesan: any;
  private startTime = {hour: 9, minute: 0};
  private endTime = {hour: 11, minute: 0};
  private namaMatkul : any[] = [];
  private jam_mulai: any;
  private menit_mulai: any;
  private jam_akhir: any;
  private menit_akhir: any;
  private waktu: any;
  private dataPesanan: any[] = [];
  private dataPesanRuangan: any[] = [];
  private dataPesanMatkul: any[] = [];

  constructor(public http: Http, public router: Router, public dataService: DataService) { }

  ngOnInit() {
  	this.dt = new Date()
  	this.dt.setHours(14)
  	this.tanggal = this.dt.toISOString()
  	this.listruanganbook(this.tanggal)
  }

  listruanganbook(tanggal){
	let sub = "T"
	if(!(tanggal.includes(sub))){
		this.tanggal = tanggal + "T14:00:00.000Z"
	}
  	this.dataMataKuliah = [];
  	this.dataruangan = [];
  	this.http.get(this.dataService.urlGetListRuanganBook+"/"+ this.tanggal, {withCredentials: true})
  	.subscribe(res=>{
  		let data = JSON.parse(res['_body']);
  		if(data!=null){
  			this.dataruangan = data['ruangan'];
  			if(data['datajadwalmatkul'] != null){
	  			for(let i = 0; i < data['datajadwalmatkul'].length; i++){
		  			const _dataMatkul = {
		              idMatkul: data['datajadwalmatkul'][i].idmatakuliah,
		              kodeMatkul: data['datajadwalmatkul'][i].kodematakuliah,
		              namaMatkul: data['datajadwalmatkul'][i].namamatakuliah,
		              waktu: data['jadwalkuliah'][i].waktu,
		              idRuangan: data['jadwalkuliah'][i].idruangan
		            }
	            	this.dataMataKuliah.push(_dataMatkul);
	          	}
  			}
  			if(data['datapesanamatkul']!=null){
	          	for(let i = 0; i < data['datapesanamatkul'].length; i++){
		  			const _dataMatkul = {
		              idMatkul: data['datapesananmatkul'][i].idmatakuliah,
		              kodeMatkul: data['datapesananmatkul'][i].kodematakuliah,
		              namaMatkul: data['datapesananmatkul'][i].namamatakuliah,
		              waktu: data['pesananruangan'][i].waktu,
		              idRuangan: data['pesananruangan'][i].idruangan
		            }
	            	this.dataMataKuliah.push(_dataMatkul);
				}   				
  			}
  		}
	});
	this.http.get(this.dataService.urlMyPesananRuangan, {withCredentials: true})
	.subscribe(res=>{
		let data = JSON.parse(res['_body']);
		if(data!=null){
			this.dataPesanan = data['datapesanan']
			this.dataPesanRuangan = data['dataruangan']
			this.dataPesanMatkul = data['datamatkul']
		}
	});
  }

  bookRuangan(){
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
	  let creds = JSON.stringify({
		idruangan: Number(this.idruangan),
		idmatakuliah: Number(this.dataMatkul[this.matkul],)
		tanggal: this.tglpesan,
		waktu: this.waktu,
		status: "Menunggu konfirmasi",
		alasan: this.alasan
	  });
	  var headers = new Headers();
	  headers.append("Content-Type", "application/json");
	  this.http.post(this.dataService. urlAddPesananRuangan, creds, {withCredentials:true,headers: headers})
  	  .subscribe(res=>{
		if(res['status'] == 201){
			swal({
				title : 'Berhasil ditambahkan',
				type : 'success'
			});
			this.dataMatkul = [];
			this.namaMatkul = [];
			this.matkul = "";
			this.tglpesan = "";
			this.waktu = "";
			this.alasan = "";
			this.ngOnInit();
		} else if(res['_body'] == "duplicate"){
			console.log("gagal");
		} 
	   });
	  
  }

  setidruangan(id){
	  this.idruangan = id
	  this.namaMatkul = []
	  this.http.get(this.dataService.urlGetJadwalKuliah, {withCredentials: true})
      .subscribe(res => {
        let data = JSON.parse(res['_body']);
        if(data != null){
          for(let i = 0; i < data['datamatakuliah'].length; i++){
            this.namaMatkul.push(data['datamatakuliah'][i].kodematakuliah + " " + data['datamatakuliah'][i].namamatakuliah + " (" + data['datajadwalkuliah'][i].KelasParalel + ")");
            this.dataMatkul[data['datamatakuliah'][i].kodematakuliah + " " + data['datamatakuliah'][i].namamatakuliah + " (" + data['datajadwalkuliah'][i].KelasParalel + ")"] = data['datamatakuliah'][i].idmatakuliah;
          }
        }

      })
  }


}
