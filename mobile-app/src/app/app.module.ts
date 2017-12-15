import { BrowserModule } from '@angular/platform-browser';
import { ErrorHandler, NgModule } from '@angular/core';
import { IonicApp, IonicErrorHandler, IonicModule } from 'ionic-angular';
import { SplashScreen } from '@ionic-native/splash-screen';
import { StatusBar } from '@ionic-native/status-bar';

import { MyApp } from './app.component';
import { AuthPage } from '../pages/auth/auth';
import { ChatroomPage } from '../pages/chatroom/chatroom';
import { HomePage } from '../pages/home/home';
import { JadwalPage } from '../pages/jadwal/jadwal';

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
    IonicModule.forRoot(MyApp)
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
    {provide: ErrorHandler, useClass: IonicErrorHandler}
  ]
})
export class AppModule {}
