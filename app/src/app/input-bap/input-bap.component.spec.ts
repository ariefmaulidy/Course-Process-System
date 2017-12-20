import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { InputBapComponent } from './input-bap.component';

describe('InputBapComponent', () => {
  let component: InputBapComponent;
  let fixture: ComponentFixture<InputBapComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ InputBapComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InputBapComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
