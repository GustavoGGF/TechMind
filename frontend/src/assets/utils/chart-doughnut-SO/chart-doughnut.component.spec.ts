import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChartDoughnutSOComponent } from './chart-doughnut.component';

describe('ChartDoughnutComponent', () => {
  let component: ChartDoughnutSOComponent;
  let fixture: ComponentFixture<ChartDoughnutSOComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ChartDoughnutSOComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ChartDoughnutSOComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
