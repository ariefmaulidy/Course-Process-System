import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router } from '@angular/router';
import { DataService } from '../data/data.service';
import { forEach } from '@angular/router/src/utils/collection';

@Component({
  selector: 'app-forumdiskusi',
  templateUrl: './forumdiskusi.component.html',
  styleUrls: ['./forumdiskusi.component.css'],
  providers: [DataService]
})
export class ForumdiskusiComponent implements OnInit {

	private halaman: number = 1;
	private listCGD: any[] = [];
  private semester: any;
  private listCGDBySemester: any[] = [];
  private checkSemua = "Semua";
  private CGDBySearch;
  private search = false;

  private checknol = 0;
	//jumlahpesan dari user belum di retrieve

  constructor(public http: Http, public router: Router, public dataService: DataService) { }


  ngOnInit() {
    this.search = false;
    this.semester = "Semua";
  	window.scrollTo(0,0);
  	this.getAllCGD();
  }

  getBySemester(sem){
    this.CGDBySearch = "";
    this.search = false;
    if (sem == "Semua"){
        this.listCGDBySemester = [];
        this.semester = "Semua";
        this.listCGD = [];
        this.getAllCGD();
      } else {
        this.listCGDBySemester = [];
        this.semester = sem;
        for(let i = 0; i < this.listCGD.length; i++){
          if(this.listCGD[i].semester == sem){
            const CGD = {
                  idCGD: this.listCGD[i].idCGD,
                  namaMatkul :  this.listCGD[i].namaMatkul,
                  jumlahPesan : this.listCGD[i].jumlahPesan,
                  kodeMatkul : this.listCGD[i].kodeMatkul,
                  semester : this.listCGD[i].semester
                }
            this.listCGDBySemester.push(CGD);
            }
          }
         
      }
  }

  getBySearch(keySearch){
    if(keySearch != ""){
      this.search = false;
      this.CGDBySearch = null;
      for(let i=0; i< this.listCGD.length; i++){
        if(this.listCGD[i].kodeMatkul.toUpperCase() == keySearch.toUpperCase()){
          const CGD = {
              idCGD: this.listCGD[i].idCGD,
              namaMatkul :  this.listCGD[i].namaMatkul,
              jumlahPesan : this.listCGD[i].jumlahPesan,
              kodeMatkul : this.listCGD[i].kodeMatkul,
              semester : this.listCGD[i].semester
             }
            this.CGDBySearch = CGD;
            this.search = true;
            break;
            } 
        }
      }
    }
  

  getAllCGD(){
    this.CGDBySearch = "";
  	this.http.get(this.dataService.urlGetCGD, {withCredentials: true})
  		.subscribe(res => {
  			let data = JSON.parse(res['_body']);
  			if(data != null){
  				for(let i = 0; i < data['datacgd'].length; i++){
  					const CGD = {
		              idCGD: data['datacgd'][i].idcgd,
		              namaMatkul :  data['datamatakuliah'][i].kodematakuliah+ " "+ data['datamatakuliah'][i].namamatakuliah + "(" + data['datajadwalkuliah'][i].KelasParalel + ")",
            		  jumlahPesan : data['datacgd'][i].jumlahpesan,
            		  kodeMatkul : data['datamatakuliah'][i].kodematakuliah,
                  semester : data['datamatakuliah'][i].semester,
                  unread : data['unread'][i]
                }
            this.listCGD.push(CGD);
          }
  			}

  		})
  }

}
