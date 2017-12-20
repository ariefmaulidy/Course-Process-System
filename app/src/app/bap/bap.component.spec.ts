import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BapComponent } from './bap.component';

describe('BapComponent', () => {
  let component: BapComponent;
  let fixture: ComponentFixture<BapComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BapComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BapComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
