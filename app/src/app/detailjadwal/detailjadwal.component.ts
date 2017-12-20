import { Component, OnInit } from '@angular/core';
import { DataService } from '../data/data.service';
import { Http, Headers } from '@angular/http';
import { Router , ActivatedRoute} from '@angular/router';


@Component({
  selector: 'app-detailjadwal',
  templateUrl: './detailjadwal.component.html',
  styleUrls: ['./detailjadwal.component.css'],
  providers:[DataService]
})
export class DetailjadwalComponent implements OnInit {
  private sub: any;
  private id_matkul: number;

  private matakuliah;
  private pengajar;
  private namaruangan;
  private nim : any[] = [];
  private nama_mahasiswa : any[] = [];
  private status : any[] = [];



  constructor(public http: Http, public route: ActivatedRoute, public router: Router,public dataService: DataService) {
  }

  ngOnInit() {
  	this.sub = this.route.params.subscribe(params => {
       this.id_matkul = +params['id_matkul'];
  	});
  }

  getDetailMatkul() {
  	this.http.get(this.dataService.urlGetJadwalKuliah + "/" + this.id_matkul, {withCredentials: true})
	    .subscribe(res =>{
	    	if(res['_body'] != "belum login"){
		    	let data = JSON.parse(res['_body']);
		    	console.log(data);
	    	}
	    });
  }

}
