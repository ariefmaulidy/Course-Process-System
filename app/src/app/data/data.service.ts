import { Injectable } from '@angular/core';
import { Subject }    from 'rxjs/Subject';

@Injectable()
export class DataService{
	public loggedIn: boolean;
	public token;
	public usernameNow: string;

	public hostname = 'http://localhost:8080/';
	public hostWebsocket = 'ws://localhost:8080/';

	//route API

	//user
	public urlLogin = this.hostname + 'login';
	public urlLogout = this.hostname + 'logout';
	public urlCheckExpiredToken = this.hostname + 'checkexpiredtoken';
	
	//bap blum fix
	public urlAddBAP = this.hostname + 'addbap';
	public urlGetBAP = this.hostname + 'bap'; //untuk all tidak ditambahkan parameter
	public urlEditBAP = this.hostname + 'editbap';
	public urlGetJadwalBAP = this.hostname + 'jadwalbap';

	//cgd
	public urlGetCGD = this.hostname + 'roomcgd'; //untuk all tidak ditambahkan parameter

	//matakuliah
	public urlGetMataKuliah = this.hostname + 'matakuliah';

	//jadwalkuliah
	public urlGetJadwalKuliah = this.hostname + 'detailjadwalkuliah'; //untuk all tidak ditambahkan parameter

	//pengajar
	public urlGetPengajar = this.hostname + 'pengajar'; //untuk all tidak ditambahkan parameter
	public urlGetJadwalMengajar = this.hostname + 'jadwalmengajar';
	public urlGetPengajarBAP = this.hostname + 'getpengajar';

	//pengelolaruangan
	public urlPersetujuanPesanan = this.hostname + 'persetujuanpesanan';
	public urlPenolakanPesanan = this.hostname + 'penolakanpesanan';

	//pesan
	public urlAddPesan = this.hostname + 'addpesan'; //parameter idcgd

	//pesananruangan yang getnya belum buat di TU
	public urlAddPesananRuangan = this.hostname + 'addpesananruangan';
	public urlMyPesananRuangan = this.hostname + 'mypesananruangan';

	//tatausaha
	public urlAssignPJKelas = this.hostname + 'assignpjkelas';
	public urlGetListRuanganBook = this.hostname + 'listbookruangan'; //paramnya tanggal
	public urlEditPJKelas = this.hostname + 'editpjkelas'; 

	//tempcgd
	public urlTempCGD = this.hostname + 'tempcgd'; //untuk all tidak ditambahkan parameter

	//socket
	public urlSocket = this.hostWebsocket + 'roomcgd'; //parameter idcgd

	public loginState(cek){
		this.loggedIn = cek;
	}

	public loginToken(cek){
		this.token = cek;
	}

	public loginUser(cek){
		this.usernameNow = cek;
	}
}