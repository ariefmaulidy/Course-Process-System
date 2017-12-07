import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ForumdiskusiComponent } from './forumdiskusi.component';

describe('ForumdiskusiComponent', () => {
  let component: ForumdiskusiComponent;
  let fixture: ComponentFixture<ForumdiskusiComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ForumdiskusiComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ForumdiskusiComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
