import { BrowserModule } from '@angular/platform-browser';
import { ErrorHandler, NgModule } from '@angular/core';
import { IonicApp, IonicErrorHandler, IonicModule } from 'ionic-angular';
import { SplashScreen } from '@ionic-native/splash-screen';
import { StatusBar } from '@ionic-native/status-bar';
import { HttpModule } from '@angular/http';
import { HttpClientModule } from '@angular/common/http';

import { MyApp } from './app.component';
import { AuthPage } from '../pages/auth/auth';
import { ChatroomPage } from '../pages/chatroom/chatroom';
import { HomePage } from '../pages/home/home';
import { JadwalPage } from '../pages/jadwal/jadwal';
import { DataProvider } from '../providers/data/data';

@NgModule({
  declarations: [
    MyApp,
    AuthPage,
    ChatroomPage,
    HomePage,
    JadwalPage
  ],
  imports: [
    BrowserModule,
    IonicModule.forRoot(MyApp),
    HttpModule,
    HttpClientModule
  ],
  bootstrap: [IonicApp],
  entryComponents: [
    MyApp,
    AuthPage,
    ChatroomPage,
    HomePage,
    JadwalPage
  ],
  providers: [
    StatusBar,
    SplashScreen,
    {provide: ErrorHandler, useClass: IonicErrorHandler},
    DataProvider
  ]
})
export class AppModule {}
