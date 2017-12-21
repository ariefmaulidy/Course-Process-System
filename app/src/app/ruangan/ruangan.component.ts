import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router } from '@angular/router';
import { DataService } from '../data/data.service';



@Component({
  selector: 'app-ruangan',
  templateUrl: './ruangan.component.html',
  styleUrls: ['./ruangan.component.css'],
  providers: [DataService]
})
export class RuanganComponent implements OnInit {

  private dt;
  private tanggal;
  private dataruangan: any[] = [];
  private dataMataKuliah: any[] = [];
  private halaman: number = 1;

  constructor(public http: Http, public router: Router, public dataService: DataService) { }

  ngOnInit() {
  	this.dt = new Date()
  	this.dt.setHours(14)
  	this.tanggal = this.dt.toISOString()
  	this.listruanganbook(this.tanggal)
  }

  listruanganbook(tanggal){
  	this.tanggal = tanggal + "T14:00:00.000Z"
  	this.dataMataKuliah = [];
  	this.dataruangan = [];
  	this.http.get(this.dataService.urlGetListRuanganBook+"/"+ this.tanggal, {withCredentials: true})
  	.subscribe(res=>{
  		let data = JSON.parse(res['_body']);
  		console.log(data);
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
        	console.log(this.dataMataKuliah)
  		}
  	});
  }

}
