import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { JadwalkuliahComponent } from './jadwalkuliah.component';

describe('JadwalkuliahComponent', () => {
  let component: JadwalkuliahComponent;
  let fixture: ComponentFixture<JadwalkuliahComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ JadwalkuliahComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(JadwalkuliahComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
