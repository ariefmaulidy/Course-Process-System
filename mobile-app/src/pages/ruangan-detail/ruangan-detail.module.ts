import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { RuanganDetailPage } from './ruangan-detail';

@NgModule({
  declarations: [
    RuanganDetailPage,
  ],
  imports: [
    IonicPageModule.forChild(RuanganDetailPage),
  ],
})
export class RuanganDetailPageModule {}
