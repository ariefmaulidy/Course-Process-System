import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DetailjadwalComponent } from './detailjadwal.component';

describe('DetailjadwalComponent', () => {
  let component: DetailjadwalComponent;
  let fixture: ComponentFixture<DetailjadwalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DetailjadwalComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DetailjadwalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
