import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DevicesDetailsComponent } from './devices-details.component';

describe('DevicesDetailsComponent', () => {
  let component: DevicesDetailsComponent;
  let fixture: ComponentFixture<DevicesDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DevicesDetailsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DevicesDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
