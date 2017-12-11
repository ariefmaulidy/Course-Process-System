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

	//cgd
	public urlGetCGD = this.hostname + 'cgd'; //untuk all tidak ditambahkan parameter

	//jadwalkuliah
	public urlGetJadwalKuliah = this.hostname + 'detailjadwalkuliah'; //untuk all tidak ditambahkan parameter

	//pengajar
	public urlGetPengajar = this.hostname + 'pengajar'; //untuk all tidak ditambahkan parameter
	public urlGetJadwalMengajar = this.hostname + 'jadwalmengajar';

	//pengelolaruangan
	public urlPersetujuanPesanan = this.hostname + 'persetujuanpesanan';
	public urlPenolakanPesanan = this.hostname + 'penolakanpesanan';

	//pesan
	public urlAddPesan = this.hostname + 'addpesan'; //parameter idcgd

	//pesananruangan yang getnya belum buat di TU
	public urlAddPesananRuangan = this.hostname + 'addpesananruangan';

	//tatausaha
	public urlAssignPJKelas = this.hostname + 'assignpjkelas';
	public urlGetListRuanganBook = this.hostname + 'listbookruangan'; //paramnya tanggal

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