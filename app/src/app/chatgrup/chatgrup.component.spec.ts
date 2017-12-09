import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ChatgrupComponent } from './chatgrup.component';

describe('ChatgrupComponent', () => {
  let component: ChatgrupComponent;
  let fixture: ComponentFixture<ChatgrupComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ChatgrupComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ChatgrupComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
