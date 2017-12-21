import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router } from '@angular/router';
import { DataService } from '../data/data.service';
import { forEach } from '@angular/router/src/utils/collection';


@Component({
  selector: 'app-jadwalkuliah',
  templateUrl: './jadwalkuliah.component.html',
  styleUrls: ['./jadwalkuliah.component.css'],
  providers: [DataService]
})
export class JadwalkuliahComponent implements OnInit {

	private halaman: number = 1;
  private dataMataKuliah: any[] = [];
  private semester: any;
  private dataMataKuliahBySemester: any[] = [];
  private checkSemua = "Semua";
  private matkulBySearch;
  private search = false;

  	constructor(public http: Http, public router: Router, public dataService: DataService) { }

  	ngOnInit() {
      this.search = false;
      this.semester = "Semua";
  		window.scrollTo(0,0);
  		this.getAllJadwalKuliah();
  	}

    getBySemester(sem){
      this.matkulBySearch = "";
      this.search = false;
      if (sem == "Semua"){
        this.dataMataKuliahBySemester = [];
        this.semester = "Semua";
        this.dataMataKuliah = [];
        this.getAllJadwalKuliah();
      } else {
        this.dataMataKuliahBySemester = [];
        this.semester = sem;
        let count = 0;
        for(let i = 0; i < this.dataMataKuliah.length; i++){
          if(this.dataMataKuliah[i].semester == sem){
            count = count + 1;
            const _dataMatkul = {
              no: count,
              idMatkul: this.dataMataKuliah[i].idMatkul,
              kodeMatkul: this.dataMataKuliah[i].kodeMatkul,
              namaMatkul: this.dataMataKuliah[i].namaMatkul,
              hari: this.dataMataKuliah[i].hari,
              waktu: this.dataMataKuliah[i].waktu,
              namaRuangan: this.dataMataKuliah[i].namaRuangan,
              semester : sem
            }
            this.dataMataKuliahBySemester.push(_dataMatkul);
            }
          }
         
      }
    }

    getBySearch(keySearch){
    if(keySearch != ""){
      this.search = false;
      this.matkulBySearch = null;
      let count = 0;
      for(let i=0; i< this.dataMataKuliah.length; i++){
        if(this.dataMataKuliah[i].kodeMatkul.toUpperCase() == keySearch.toUpperCase()){
          count = count + 1;
            const _dataMatkul = {
              no: count,
              idMatkul: this.dataMataKuliah[i].idMatkul,
              kodeMatkul: this.dataMataKuliah[i].kodeMatkul,
              namaMatkul: this.dataMataKuliah[i].namaMatkul,
              hari: this.dataMataKuliah[i].hari,
              waktu: this.dataMataKuliah[i].waktu,
              namaRuangan: this.dataMataKuliah[i].namaRuangan,
              semester : this.dataMataKuliah[i].semester
            }
            this.matkulBySearch = _dataMatkul;
            this.search = true;
            break;
            } 
        }
      }
    }

  	getAllJadwalKuliah(){
      this.matkulBySearch = "";
  		this.http.get(this.dataService.urlGetJadwalKuliah, {withCredentials: true})
  		.subscribe(res => {
  			let data = JSON.parse(res['_body']);
  			if(data != null){
  				for(let i = 0; i < data['datamatakuliah'].length; i++){
  					const _dataMatkul = {
              no: i+1,
              idMatkul: data['datamatakuliah'][i].idmatakuliah,
              kodeMatkul: data['datamatakuliah'][i].kodematakuliah,
              namaMatkul: data['datamatakuliah'][i].namamatakuliah + " (" + data['datajadwalkuliah'][i].KelasParalel + ")",
              hari: data['datajadwalkuliah'][i].hari,
              waktu: data['datajadwalkuliah'][i].waktu,
              namaRuangan: data['dataruangan'][i].namaruangan,
              semester : data['datamatakuliah'][i].semester
            }
            this.dataMataKuliah.push(_dataMatkul);
          }
  			}

  		})
  	}

}
