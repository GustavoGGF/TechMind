import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChatPointGetMachinesDayComponent } from './chat-point-get-machines-day.component';

describe('ChatPointGetMachinesDayComponent', () => {
  let component: ChatPointGetMachinesDayComponent;
  let fixture: ComponentFixture<ChatPointGetMachinesDayComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ChatPointGetMachinesDayComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ChatPointGetMachinesDayComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
