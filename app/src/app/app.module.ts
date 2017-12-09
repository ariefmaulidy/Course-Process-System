import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { NavbarComponent } from './navbar/navbar.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { JadwalkuliahComponent } from './jadwalkuliah/jadwalkuliah.component';
import { AuthComponent } from './auth/auth.component';
import { DetailjadwalComponent } from './detailjadwal/detailjadwal.component';
import { ForumdiskusiComponent } from './forumdiskusi/forumdiskusi.component';
import { RuanganComponent } from './ruangan/ruangan.component';
import { ChatgrupComponent } from './chatgrup/chatgrup.component';

@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    DashboardComponent,
    JadwalkuliahComponent,
    AuthComponent,
    DetailjadwalComponent,
    ForumdiskusiComponent,
    RuanganComponent,
    ChatgrupComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    RouterModule.forRoot([
      {
        path :'',
        component:DashboardComponent
      },    
      {
        path: 'dashboard',
        component: DashboardComponent
      },
      {
        path : 'jadwalkuliah',
        component: JadwalkuliahComponent
      },
      {
        path : 'detailjadwal',
        component: DetailjadwalComponent
      },
      {
        path : 'forumdiskusi',
        component: ForumdiskusiComponent
      },
      {
        path : 'ruangan',
        component: RuanganComponent
      },
      {
        path : 'chatgrup',
        component: ChatgrupComponent
      },
      {
        path : 'auth',
        component : AuthComponent
      }
           
    ], { useHash: true })

  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
