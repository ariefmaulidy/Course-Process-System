import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Router, ActivatedRoute } from '@angular/router';
import { DataService } from '../data/data.service';
import { forEach } from '@angular/router/src/utils/collection';
import swal from 'sweetalert2'

@Component({
  	selector: 'app-detailjadwal',
  	templateUrl: './detailjadwal.component.html',
  	styleUrls: ['./detailjadwal.component.css'],
	providers: [DataService]
})
export class DetailjadwalComponent implements OnInit {

	private idjadwalkuliah;
	private mataKuliah;
	private pengajar;
	private namaRuangan;
	private SKS;
	private kelas;

	public id:number;
	public sub: any;

	private datamahasiswa: any[] = [];

	private hasPJ;
	private idPJ;

	constructor(public http: Http, public route: ActivatedRoute, public router: Router, public dataService: DataService) { }

	ngOnInit() {
		this.hasPJ = false;
		window.scrollTo(0,0);
	  	this.sub = this.route.params.subscribe(params => {
	       this.id = +params['idmatkul'];
	  	});
	  	this.getDetailJadwal();
	}

	getDetailJadwal(){
		this.http.get(this.dataService.urlGetJadwalKuliah + '/' + this.id, {withCredentials: true})
  		.subscribe(res => {
  			let data = JSON.parse(res['_body']);
  			/*console.log(data)*/
  			if(data != null){
  				this.idjadwalkuliah = data['datajadwalkuliah'].idjadwalkuliah
  				this.mataKuliah = data['matakuliah'].namamatakuliah
  				this.SKS = data['matakuliah'].sks
  				this.kelas = data['datajadwalkuliah'].KelasParalel
  				this.pengajar = data['datadosen'].nama
				this.namaRuangan = data['dataruangan'].namaruangan
				if(data['mahasiswa'] != null){
					for(let i = 0; i < data['mahasiswa'].length; i++){
						const _dataMahasiswa = {
							nim: data['mahasiswa'][i].nim,
							nama: data['mahasiswa'][i].nama,
							statusmayor: data['datapesertakuliah'][i].statusmayor,
							statuspj: data['datapesertakuliah'][i].statuspj,
							iduser: data['mahasiswa'][i].iduser
						}
						if(_dataMahasiswa.statuspj == "Penanggung Jawab" && this.hasPJ == false){
							this.hasPJ = true;
							this.idPJ = _dataMahasiswa.iduser;
						}
						this.datamahasiswa.push(_dataMahasiswa)
					}
				}  
  			}
		});
	}

	assignPJ(idjadwalkuliah, iduser){
		this.http.put(this.dataService.urlAssignPJKelas + '/' + idjadwalkuliah + "?iduser=" + iduser ,null, {withCredentials:true})
		.subscribe(res =>{
			if (res['status'] == 204){
				swal({
				  title : 'Assigned',
				  type : 'success'
				})
				this.datamahasiswa = []
				this.ngOnInit()
			} else {
				console.log("gagal")
			}
		})
	}

	editPJ(idjadwalkuliah){
		this.http.put(this.dataService.urlEditPJKelas+ '/' + idjadwalkuliah + "?iduser=" + this.idPJ ,null, {withCredentials:true})
		.subscribe(res =>{
			if (res['status'] == 204){
				swal({
				  title : 'Pemilihan PJ diulang',
				  type : 'success'
				})
				this.datamahasiswa = []
				this.ngOnInit()
			} else {
				console.log("gagal")
			}
		})
	}
}
